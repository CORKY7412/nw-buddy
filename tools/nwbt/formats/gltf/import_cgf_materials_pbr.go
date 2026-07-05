package gltf

import (
	"fmt"
	"image/color"
	"log/slog"
	"math"

	"nw-buddy/tools/formats/gltf/extensions"
	"nw-buddy/tools/formats/image"
	"nw-buddy/tools/formats/mtl"
	"nw-buddy/tools/utils"
	"path"
	"slices"
	"strings"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/ext/texturetransform"
)

// Tries to interpret original material intent as best as possible
// and map it to a PBR material or whatever is closest to it in glTF.
//
// This works well for most new world materials, but does not cover all shader features.
// Good enough for a preview
func (doc *Document) augmentMaterialWithPbr(material *gltf.Material, cgfMat *mtl.Material, options ImportCgfMaterialsOptions) {
	diffuseOnly := options.DiffuseOnly
	textureBaking := options.TextureBaking

	if material.Extensions == nil {
		material.Extensions = gltf.Extensions{}
	}

	material.AlphaMode = gltf.AlphaOpaque
	material.AlphaCutoff = gltf.Float(0.5)
	material.DoubleSided = false
	material.PBRMetallicRoughness = &gltf.PBRMetallicRoughness{
		BaseColorFactor: float4(1, 1, 1, 1),
		MetallicFactor:  gltf.Float(0), // non metallic flow for NW
		RoughnessFactor: gltf.Float(1),
	}
	material.Extensions[extensions.KHR_materials_ior] = extensions.KHRMaterialsIOR{
		IOR: gltf.Float(float64(options.CustomIOR)), // specular workflow requires this
	}
	doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, extensions.KHR_materials_ior)

	// material.Extras = ExtrasStore(material.Extras, "mtl", cgfMat)

	mtlDiffuse := cgfMat.Diffuse
	mtlSpecular := cgfMat.Specular
	// mtlEmissive := mtl.ParamColor(m.Emissive) // HINT: Emissive is not used in NW
	mtlEmittance := cgfMat.Emittance
	for i := range mtlEmittance {
		// WTF:
		// Emittance="-341776551831207702696034304.000000,0.000000,0.000000,0.000000"
		if mtlEmittance[i] < 0 {
			mtlEmittance[i] = 0
		}
	}

	mtlOpacity := float32(1)
	if cgfMat.Opacity != nil {
		mtlOpacity = *cgfMat.Opacity
	}
	mtlAlphaTest := float32(1)
	if cgfMat.AlphaTest != nil {
		mtlAlphaTest = *cgfMat.AlphaTest
	}
	mtlShininess := float32(255)
	if cgfMat.Shininess != nil {
		mtlShininess = *cgfMat.Shininess
	}

	mapDiffuse := cgfMat.TextureByMapType(mtl.MtlMap_Diffuse)
	mapBumpmap := cgfMat.TextureByMapType(mtl.MtlMap_Bumpmap)
	mapSmoothness := cgfMat.TextureByMapType(mtl.MtlMap_Smoothness)
	mapSpecular := cgfMat.TextureByMapType(mtl.MtlMap_Specular)
	mapCustom := cgfMat.TextureByMapType(mtl.MtlMap_Custom)
	mapEmittance := cgfMat.TextureByMapType(mtl.MtlMap_Emittance)
	// mapOpacity := m.TextureWithMap(mtl.MtlMap_Opacity)
	mapDecal := cgfMat.TextureByMapType(mtl.MtlMap_Decal)

	if diffuseOnly {
		mapBumpmap = nil
		mapSmoothness = nil
		mapSpecular = nil
		mapCustom = nil
		mapEmittance = nil
		mapDecal = nil
	}

	shader := strings.ToLower(cgfMat.Shader)
	isGlass := shader == "glass"
	isHair := shader == "hair"
	isFxTransp := shader == "fxmeshadvancedtransp" // only isabella gaze has this shader
	isGeometryFog := shader == "geometryfog"
	isGeometryBeam := shader == "geometrybeam" || shader == "geometrybeamsimple"
	isParticle := shader == "meshparticle"
	isNoDraw := shader == "nodraw"
	isVegetation := shader == "vegetation"

	features := strings.Split(cgfMat.StringGenMask, "%")
	shaderOverlayMask := slices.Contains(features, "OVERLAY_MASK")
	shaderEmittanceMap := slices.Contains(features, "EMITTANCE_MAP")
	shaderEmissiveDecal := slices.Contains(features, "EMISSIVE_DECAL")
	shaderTintColorMap := slices.Contains(features, "TINT_COLOR_MAP")
	shaderVertColors := slices.Contains(features, "VERTCOLORS")
	shaderDecal := slices.Contains(features, "DECAL")
	if isVegetation {
		// vegetation shader uses vert colors differently
		shaderVertColors = false
	}

	if shaderDecal && cgfMat.Params != nil {
		// PublicParams EmittanceMapGamma="1" DecalDiffuseOpacity="0" DecalFalloff="1" DecalAlphaMult="0" IndirectColor="0.24620135,0.24620135,0.24620135,0.25"/>
		// decalOpacity := float32(1) // DecalDiffuseOpacity -> unused in shader
		// decalFalloff := float32(1) // DecalFalloff
		// decalAlphaMult := float32(1) // DecalAlphaMult
		// if v, ok := m.Params.LoadFloat("DecalFalloff"); ok {
		// 	decalFalloff = v
		// }
		// if v, ok := m.Params.LoadFloat("DecalAlphaMult"); ok {
		// 	decalAlphaMult = v
		// }

		// does break other stuff, so disable for now
		// mtlOpacity = float32(math.Pow(float64(mtlOpacity*decalAlphaMult), float64(decalFalloff)))
	}

	// HINT: Bumpmap and Smoothness are usually the same texture source
	if mapSmoothness == nil && mapBumpmap != nil && strings.Contains(mapBumpmap.File, "_ddna") {
		mapSmoothness = mapBumpmap
	}
	if isFxTransp && mapDiffuse == nil {
		mapDiffuse = mapCustom
		mapCustom = nil
	}
	if shaderTintColorMap && (mapDiffuse == nil || strings.Contains(mapDiffuse.File, "/white.dds")) {
		// replaces EngineAssets/Textures/white.dds texture
		mapDiffuse = mapCustom
		mapCustom = nil
	}
	if shaderEmissiveDecal && mapEmittance == nil {
		mapEmittance = mapDecal
		mapDecal = nil
	}

	// texMods := make(map[string]any)
	// addTexMod := func(tex *gltf.Texture, mTex *mtl.Texture) {
	// 	if tex == nil || mTex == nil || mTex.TexMod == nil {
	// 		return
	// 	}
	// 	refId, _ := ExtrasLoad[string](tex.Extras, ExtraKeyRefID)
	// 	texMods[refId] = mTex.TexMod
	// 	material.Extras = ExtrasStore(material.Extras, "mods", texMods)
	// }

	var texCustom *gltf.Texture
	// var texCustomTransform *texturetransform.TextureTranform
	if mapCustom != nil {
		texCustom = doc.LoadOrStoreMtlTexture(mapCustom)
		// addTexMod(texCustom, mapCustom)
		// texCustomTransform = resolveTextureTransform(mapCustom)
	}
	var texDiffuse *gltf.Texture
	var texDiffuseTransform *texturetransform.TextureTranform
	if mapDiffuse != nil {
		texDiffuse = doc.LoadOrStoreMtlTexture(mapDiffuse)
		texDiffuseTransform = resolveTextureTransform(mapDiffuse)
		// addTexMod(texDiffuse, mapDiffuse)
	}
	var texBumpmap *gltf.Texture
	var texBumpmapTransform *texturetransform.TextureTranform
	if mapBumpmap != nil {
		texBumpmap = doc.LoadOrStoreMtlTexture(mapBumpmap)
		texBumpmapTransform = resolveTextureTransform(mapBumpmap)
		if texBumpmapTransform == nil {
			texBumpmapTransform = texDiffuseTransform
		}
		// addTexMod(texBumpmap, mapBumpmap)
	}

	var texSpecular *gltf.Texture
	var texSpecularTransform *texturetransform.TextureTranform
	if mapSpecular != nil {
		texSpecular = doc.LoadOrStoreMtlTexture(mapSpecular)
		texSpecularTransform = resolveTextureTransform(mapSpecular)
		if texBumpmapTransform == nil {
			texSpecularTransform = texDiffuseTransform
		}
		// addTexMod(texSpecular, mapSpecular)
	}
	texTransformPublic := resolveTextureTransformPublic(cgfMat.Params)

	var texSmoothnessTransform *texturetransform.TextureTranform
	if mapSmoothness != nil {
		texSmoothnessTransform = resolveTextureTransform(mapSmoothness)
		if texSmoothnessTransform == nil {
			texSmoothnessTransform = texDiffuseTransform
		}
	}
	if textureBaking && texSpecular != nil && mapSmoothness != nil {
		key := buildTexturePath(texSpecular, mapSmoothness.File, mapSmoothness.AssetId)
		tex, err := doc.LoadOrStoreTexture(nil, key, func() ([]byte, error) {
			specBytes, err := doc.LoadTextureBytes(texSpecular)
			if err != nil {
				return nil, err
			}
			if specBytes == nil {
				return nil, fmt.Errorf("specBytes is nil")
			}
			smoothImg, err := doc.ResolveAndLoadImage(mapSmoothness)
			if err != nil {
				return nil, err
			}
			if smoothImg == nil {
				slog.Warn("smoothImg is nil", "file", key)
				return specBytes, nil
			}
			if smoothImg.Alpha == nil {
				slog.Warn("smoothImg.Alpha is nil", "file", key)
				return specBytes, nil
			}
			return image.MixPngImages(specBytes, smoothImg.Alpha, func(spec, smooth color.NRGBA) color.NRGBA {
				return color.NRGBA{
					spec.R,
					spec.G,
					spec.B,
					smooth.R, // HINT: don't read the alpha channel. Its a grayscale image, hence just access .R
				}
			})
		})
		if err != nil {
			slog.Warn("Can't combine spec+smooth", "file", key, "err", err)
		} else {
			texSpecular = tex
		}
	}

	var texEmissive *gltf.Texture
	if textureBaking && shaderEmissiveDecal && mapEmittance != nil && mapDecal != nil {
		key := mapEmittance.File + "_" + hashString(path.Join(mapDecal.File, mapDecal.AssetId))
		tex, err := doc.LoadOrStoreTexture(nil, key, func() ([]byte, error) {
			emitImg, err := doc.ResolveAndLoadImage(mapEmittance)
			if err != nil {
				return nil, err
			}
			decalImg, err := doc.ResolveAndLoadImage(mapDecal)
			if err != nil {
				return nil, err
			}
			return image.MixPngImages(emitImg.Data, decalImg.Data, func(emit, decal color.NRGBA) color.NRGBA {
				return color.NRGBA{
					uint8(float32(emit.R) / 255.0 * float32(decal.R)),
					uint8(float32(emit.G) / 255.0 * float32(decal.G)),
					uint8(float32(emit.B) / 255.0 * float32(decal.B)),
					uint8(float32(emit.A) / 255.0 * float32(decal.A)),
				}
			})
		})
		if err != nil {
			slog.Warn("Can't combine emit+decal", "file", key, "err", err)
		} else {
			texEmissive = tex
		}
	}

	if (shaderEmissiveDecal || shaderEmittanceMap) && texEmissive == nil && mapEmittance != nil {
		texEmissive = doc.LoadOrStoreMtlTexture(mapEmittance)
	}

	if len(mtlDiffuse) >= 3 {
		factor := [4]float64{
			float64(mtlDiffuse[0]),
			float64(mtlDiffuse[1]),
			float64(mtlDiffuse[2]),
			1,
		}
		if len(mtlDiffuse) == 4 {
			factor[3] = float64(mtlDiffuse[3])
		}
		material.PBRMetallicRoughness.BaseColorFactor = &factor
	}
	if texDiffuse != nil {
		texExt := gltf.Extensions{}
		if texDiffuseTransform != nil {
			texExt[texturetransform.ExtensionName] = &texDiffuseTransform
		} else if texTransformPublic != nil {
			texExt[texturetransform.ExtensionName] = &texTransformPublic
		}
		material.PBRMetallicRoughness.BaseColorTexture = &gltf.TextureInfo{
			Index:      slices.Index(doc.Textures, texDiffuse),
			Extensions: texExt,
		}
	}
	if texBumpmap != nil {
		texExt := gltf.Extensions{}
		if texBumpmapTransform != nil {
			texExt[texturetransform.ExtensionName] = &texBumpmapTransform
		} else if texTransformPublic != nil {
			texExt[texturetransform.ExtensionName] = &texTransformPublic
		}
		material.NormalTexture = &gltf.NormalTexture{
			Index:      gltf.Index(slices.Index(doc.Textures, texBumpmap)),
			Scale:      gltf.Float(1),
			Extensions: texExt,
		}
	}
	if !isGlass && !isHair {

		specExt := extensions.KHRMaterialsSpecular{
			SpecularFactor:       gltf.Float(1),
			SpecularTexture:      nil,
			SpecularColorFactor:  float3(1, 1, 1),
			SpecularColorTexture: nil,
		}

		if mtlShininess >= 0 {
			smoothness := float64(mtlShininess / 255)
			roughness := (1.0 - smoothness)
			specExt.SpecularFactor = gltf.Float(smoothness)
			material.PBRMetallicRoughness.RoughnessFactor = gltf.Float(roughness)
		}
		if mtlShininess == 255 && mapSmoothness != nil {
			// some materials have shininess 255 but should be rough (wood)
			// others render OK. use a compromise value of 0.5 for now
			specExt.SpecularFactor = gltf.Float(0.5)
			material.PBRMetallicRoughness.RoughnessFactor = gltf.Float(0.5)
		}
		if diffuseOnly {
			specExt.SpecularFactor = gltf.Float(0)
			material.PBRMetallicRoughness.RoughnessFactor = gltf.Float(1)
			material.PBRMetallicRoughness.MetallicFactor = gltf.Float(1)
		}

		if texSpecular != nil {
			texExt := gltf.Extensions{}
			if texSpecularTransform != nil {
				texExt[texturetransform.ExtensionName] = &texSpecularTransform
			} else if texTransformPublic != nil {
				texExt[texturetransform.ExtensionName] = &texTransformPublic
			}
			specExt.SpecularColorTexture = &gltf.TextureInfo{
				Index:      slices.Index(doc.Textures, texSpecular),
				Extensions: texExt,
			}
			specExt.SpecularTexture = &gltf.TextureInfo{
				Index:      slices.Index(doc.Textures, texSpecular),
				Extensions: texExt,
			}
		}
		if len(mtlSpecular) >= 3 {
			specExt.SpecularColorFactor = float3(float64(mtlSpecular[0]), float64(mtlSpecular[1]), float64(mtlSpecular[2]))
		}

		material.Extensions[extensions.KHR_materials_specular] = &specExt
		doc.ExtensionsRequired = utils.AppendUniq(doc.ExtensionsRequired, extensions.KHR_materials_specular)
		doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, extensions.KHR_materials_specular)
	}

	if texEmissive != nil {
		material.EmissiveFactor = [3]float64{1, 1, 1}
		material.EmissiveTexture = &gltf.TextureInfo{
			Index: slices.Index(doc.Textures, texEmissive),
		}
	}

	// (isFxTransp || texEmissive != nil) && !shaderEmissiveDecal &&
	if len(mtlEmittance) >= 4 {
		scale := math.Log(1+float64(mtlEmittance[3])) / 10
		material.EmissiveFactor = [3]float64{
			float64(mtlEmittance[0]) * scale,
			float64(mtlEmittance[1]) * scale,
			float64(mtlEmittance[2]) * scale,
		}
	} else if len(mtlEmittance) >= 3 {
		material.EmissiveFactor = [3]float64{
			float64(mtlEmittance[0]),
			float64(mtlEmittance[1]),
			float64(mtlEmittance[2]),
		}
	}

	material.AlphaMode = gltf.AlphaOpaque
	if mtlAlphaTest < 1.0 {
		cutoff := float64(mtlAlphaTest)
		material.AlphaCutoff = &cutoff
		material.AlphaMode = gltf.AlphaMask
		material.DoubleSided = true
	}
	if isGeometryFog || isGeometryBeam || isParticle || isNoDraw {
		mtlOpacity = 0.1
	}
	if mtlOpacity < 1 {
		opacity := float64(mtlOpacity)
		if material.PBRMetallicRoughness == nil {
			material.PBRMetallicRoughness = &gltf.PBRMetallicRoughness{}
		}
		if material.PBRMetallicRoughness.BaseColorFactor == nil {
			material.PBRMetallicRoughness.BaseColorFactor = &[4]float64{1, 1, 1, 1}
		}
		material.PBRMetallicRoughness.BaseColorFactor[3] = opacity
		if mtlAlphaTest == 1.0 {
			material.AlphaCutoff = nil
			material.AlphaMode = gltf.AlphaBlend
			material.DoubleSided = true
		}
	}

	if isGlass {
		factor := mtlOpacity
		if factor == 0 {
			factor = 1
		}
		material.DoubleSided = true
		material.Extensions[extensions.KHR_materials_transmission] = extensions.KHRMaterialsTransmission{
			Factor: float64(factor),
		}
		doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, extensions.KHR_materials_transmission)
	}

	if !textureBaking || (shaderOverlayMask && texCustom != nil) {
		doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, extensions.EXT_nw_material)

		ext := extensions.ExtNewWorld{
			Params:     extensions.AppearanceFromMtl(cgfMat),
			VertColors: shaderVertColors,
		}
		if texCustom != nil {
			ext.MaskTexture = &gltf.TextureInfo{
				Index: slices.Index(doc.Textures, texCustom),
			}
		}
		if !textureBaking && mapSmoothness != nil {
			texSmoothness := doc.LoadOrStoreMtlTexture(mapSmoothness, true)
			if texSmoothness != nil {
				texExt := gltf.Extensions{}
				if texSmoothnessTransform != nil {
					texExt[texturetransform.ExtensionName] = &texSmoothnessTransform
				} else if texTransformPublic != nil {
					texExt[texturetransform.ExtensionName] = &texTransformPublic
				}
				ext.SmoothTexture = &gltf.TextureInfo{
					Index:      slices.Index(doc.Textures, texSmoothness),
					Extensions: texExt,
				}
			} else {
				slog.Warn("Smoothness texture not loaded", "file", mapSmoothness.File)
			}
		}
		material.Extensions[extensions.EXT_nw_material] = ext
	}
	doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, texturetransform.ExtensionName)
}
