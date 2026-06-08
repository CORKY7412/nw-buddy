package inspect

import (
	"fmt"
	"io"
	"maps"
	"nw-buddy/tools/formats/mtl"
	"nw-buddy/tools/game"
	"nw-buddy/tools/nwfs"
	"slices"
	"strings"
	"text/tabwriter"
)

type MtlInspector struct {
	count  int
	shader map[string]*MtlShaderInfo
}

type MtlShaderInfo struct {
	count      int
	maskCount  map[string]int
	paramCount map[string]int
	paramElems map[string]int
	mapCount   map[string]int
	modCount   map[string]int
}

func NewMtlInspector() *MtlInspector {
	return &MtlInspector{
		shader: make(map[string]*MtlShaderInfo),
	}
}

func (it *MtlInspector) Inspect(assets *game.Assets, file nwfs.File) {
	it.count += 1
	f, err := mtl.Load(file)
	if err != nil {
		return
	}

	for _, m := range f.Collection() {
		shader := strings.ToLower(m.Shader)
		if it.shader[shader] == nil {
			it.shader[shader] = &MtlShaderInfo{
				maskCount:  make(map[string]int),
				paramCount: make(map[string]int),
				paramElems: make(map[string]int),
				mapCount:   make(map[string]int),
				modCount:   make(map[string]int),
			}
		}
		info := it.shader[shader]
		info.count += 1

		masks := strings.Split(m.StringGenMask, "%")
		masks = slices.DeleteFunc(masks, func(it string) bool {
			return it == ""
		})
		slices.Sort(masks)

		for _, mask := range masks {
			info.maskCount[mask] += 1
		}

		params := m.PublicParamsMap()
		paramKeys := slices.Collect(maps.Keys(params))
		slices.Sort(paramKeys)
		for _, paramKey := range paramKeys {
			info.paramCount[paramKey] += 1
			value := m.Params.Get(paramKey)
			info.paramElems[paramKey] = max(info.paramElems[paramKey], strings.Count(value, ",")+1)
		}

		for texture := range m.IterTextures() {
			info.mapCount[string(texture.Map)] += 1
			if texture.TexMod != nil {
				info.modCount[string(texture.Map)] += 1
			}
		}
	}

}

func (it *MtlInspector) Print(w io.Writer) {

	fmt.Fprintf(w, "MTL\t%d\n", it.count)

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	shader := slices.Collect(maps.Keys(it.shader))
	slices.Sort(shader)

	for _, shader := range shader {
		info := it.shader[shader]

		fmt.Fprintf(tw, "%s\t%d\n", shader, info.count)
		masks := slices.Collect(maps.Keys(info.maskCount))
		slices.SortFunc(masks, func(a, b string) int {
			return info.maskCount[b] - info.maskCount[a]
		})
		if len(masks) > 0 {
			fmt.Fprintf(tw, "  -- masks (%d) -- \n", len(masks))
			for _, mask := range masks {
				fmt.Fprintf(tw, "  %s\t%d\n", mask, info.maskCount[mask])
			}
		}

		params := slices.Collect(maps.Keys(info.paramCount))
		slices.SortFunc(params, func(a, b string) int {
			return info.paramCount[b] - info.paramCount[a]
		})
		if len(params) > 0 {
			fmt.Fprintf(tw, "  -- public params (%d) -- \n", len(params))
			for _, param := range params {
				fmt.Fprintf(tw, "  %s\t%d\t%d\n", param, info.paramCount[param], info.paramElems[param])
			}
		}

		mapNames := slices.Collect(maps.Keys(info.mapCount))
		slices.SortFunc(mapNames, func(a, b string) int {
			return info.mapCount[b] - info.mapCount[a]
		})
		if len(mapNames) > 0 {
			fmt.Fprintf(tw, "  -- maps (%d) -- \n", len(mapNames))
			for _, mapName := range mapNames {
				fmt.Fprintf(tw, "  %s\t%d\t%d\n", mapName, info.mapCount[mapName], info.modCount[mapName])
			}
		}

		fmt.Fprintf(tw, "\n")
	}
	tw.Flush()

}
