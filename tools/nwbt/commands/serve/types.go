package serve

import (
	"nw-buddy/tools/formats/catalog"
	"nw-buddy/tools/game/level"
)

type ServeApi struct {
	List                          ServeListResult          `json:"'/list'"`
	ListPattern                   ServeListResult          `json:"'/list/{filePattern}'"`
	CatalogAsset                  ServeCatalogAssetResult  `json:"'/catalog/asset/{assetId}'"`
	Stat                          ServeStatResult          `json:"'/stat/{filePattern}'"`
	LevelList                     level.LevelIndex         `json:"'/levels/list.json'"`
	LevelCoatlicueInfo            level.CoatlicueInfo      `json:"'/levels/{coatlicue}/info.json'"`
	LevelCoatlicueRegionInfo      level.RegionInfo         `json:"'/levels/{coatlicue}/{region}/info.json'"`
	LevelCoatlicueRegionCapitals  level.RegionCapitalsData `json:"'/levels/{coatlicue}/{region}/capitals.json'"`
	LevelCoatlicueRegionHeightmap []byte                   `json:"'/levels/{coatlicue}/{region}/heightmap.r16'" ts_type:"ArrayBuffer"`
	LevelCoatlicueRegionWatermap  []byte                   `json:"'/levels/{coatlicue}/{region}/watermap.r16'" ts_type:"ArrayBuffer"`
}

type ServeListResult struct {
	Items []string `json:"items"`
}

type ServeCatalogAssetResult struct {
	Asset  *catalog.Asset   `json:"asset"`
	Assets []*catalog.Asset `json:"assets"`
	Link   catalog.AssetId  `json:"link"`
	Legacy catalog.AssetId  `json:"legacy"`
}

type ServeStatResult struct {
	Items []ServeStatResultEntry `json:"items"`
}

type ServeStatResultEntry struct {
	File  string         `json:"file,omitempty"`
	Asset *catalog.Asset `json:"asset,omitempty"`
	Dds   map[string]any `json:"dds,omitempty"`
}
