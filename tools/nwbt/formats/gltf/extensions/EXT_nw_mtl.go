package extensions

import (
	"nw-buddy/tools/formats/mtl"

	"github.com/qmuntal/gltf"
)

const (
	EXT_nw_mtl = "EXT_nw_mtl" // new world material data
	EXT_nw_tex = "EXT_nw_tex" // new world texture attrubytes
)

type ExtNwMtl struct {
	Attributes mtl.MaterialAttrs   `json:"attrs,omitempty"`
	Textures   []*gltf.TextureInfo `json:"textures,omitempty"`
	Params     *mtl.PublicParams   `json:"params,omitempty"`
}
