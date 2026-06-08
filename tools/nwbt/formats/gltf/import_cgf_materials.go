package gltf

import (
	"nw-buddy/tools/formats/mtl"
	"path"

	"github.com/qmuntal/gltf"
)

type ImportCgfMaterialsOptions struct {
	TextureBaking bool
	DiffuseOnly   bool
	CustomIOR     float32
	Original      bool
}

func WithTextureBaking(value bool) func(*ImportCgfMaterialsOptions) {
	return func(o *ImportCgfMaterialsOptions) {
		o.TextureBaking = value
	}
}

func WithDiffuseOnly(value bool) func(*ImportCgfMaterialsOptions) {
	return func(o *ImportCgfMaterialsOptions) {
		o.DiffuseOnly = value
	}
}

func WithCustomIOR(value float32) func(*ImportCgfMaterialsOptions) {
	return func(o *ImportCgfMaterialsOptions) {
		o.CustomIOR = value
	}
}

func WithOriginal(value bool) func(*ImportCgfMaterialsOptions) {
	return func(o *ImportCgfMaterialsOptions) {
		o.Original = value
	}
}

func (doc *Document) ProcessMaterials(opts ...func(*ImportCgfMaterialsOptions)) {
	options := ImportCgfMaterialsOptions{
		TextureBaking: false,
		CustomIOR:     0.0,
	}
	for _, o := range opts {
		o(&options)
	}

	for _, gltfMtl := range doc.Materials {
		if !doc.IsMaterialReferenced(gltfMtl) && doc.Options.SkipUnlinkedMtl {
			continue
		}

		cgfMtl := pluckMtl(gltfMtl)
		if cgfMtl == nil {
			continue
		}

		doc.augmentMaterialWithPbr(gltfMtl, cgfMtl, options)
		doc.augmentMaterialWithOriginal(gltfMtl, cgfMtl, options)
	}
}

func lookupMtl(material *gltf.Material) *mtl.Material {
	if material == nil {
		return nil
	}
	if mtl, ok := ExtrasLoad[mtl.Material](material.Extras, "mtl"); ok {
		return &mtl
	}
	return nil
}

func pluckMtl(material *gltf.Material) *mtl.Material {
	if mtl, ok := ExtrasLoad[mtl.Material](material.Extras, "mtl"); ok {
		material.Extras = ExtrasDelete(material.Extras, "mtl")
		return &mtl
	}
	return nil
}

func float4(r, g, b, a float64) *[4]float64 {
	res := [4]float64{r, g, b, a}
	return &res
}

func float3(r, g, b float64) *[3]float64 {
	res := [3]float64{r, g, b}
	return &res
}

func buildTexturePath(tex *gltf.Texture, affix ...string) string {
	result, ok := ExtrasLoad[string](tex.Extras, ExtraKeySource)
	if !ok {
		result, ok = ExtrasLoad[string](tex.Extras, ExtraKeyRefID)
	}
	if !ok {
		result = tex.Name
	}
	return result + "_" + hashString(path.Join(affix...))
}
