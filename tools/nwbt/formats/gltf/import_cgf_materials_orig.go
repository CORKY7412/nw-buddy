package gltf

import (
	"nw-buddy/tools/formats/gltf/extensions"
	"nw-buddy/tools/formats/mtl"
	"nw-buddy/tools/utils"
	"slices"

	"github.com/qmuntal/gltf"
)

// simply translates the xml materials into a gltf structure by keeping all parameters.
func (doc *Document) augmentMaterialWithOriginal(gltfMtl *gltf.Material, cgfMat *mtl.Material, options ImportCgfMaterialsOptions) {
	nwMtl := extensions.ExtNwMtl{}
	nwMtl.Attributes = cgfMat.MaterialAttrs
	nwMtl.Params = cgfMat.Params

	for _, nwTex := range cgfMat.Textures.Texture {
		gltfTex := doc.LoadOrStoreMtlTexture(&nwTex)
		gltfTexExt := gltf.Extensions{}
		gltfTexExt[extensions.EXT_nw_tex] = nwTex
		nwMtl.Textures = append(nwMtl.Textures, &gltf.TextureInfo{
			Index:      slices.Index(doc.Textures, gltfTex),
			Extensions: gltfTexExt,
		})
	}
	if gltfMtl.Extensions == nil {
		gltfMtl.Extensions = gltf.Extensions{}
	}
	gltfMtl.Extensions[extensions.EXT_nw_mtl] = nwMtl
	doc.ExtensionsUsed = utils.AppendUniq(doc.ExtensionsUsed, extensions.EXT_nw_mtl, extensions.EXT_nw_tex)
}
