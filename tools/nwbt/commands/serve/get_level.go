package serve

import (
	"net/http"
	"nw-buddy/tools/game"
	"nw-buddy/tools/game/level"
	"nw-buddy/tools/nwfs"

	"github.com/gorilla/mux"
)

func LevelsRouter(r *mux.Router, assets *game.Assets) {
	r.HandleFunc("/list.json", getLevelsListFunc(assets))
	r.HandleFunc("/{coatlicue}/info.json", getLevelInfoFunc(assets))
	r.HandleFunc("/{coatlicue}/{region}/info.json", getRegionInfoFunc(assets))
	r.HandleFunc("/{coatlicue}/{region}/capitals.json", getRegionCapitalsFunc(assets))
	r.HandleFunc("/{coatlicue}/{region}/heightmap.r16", getRegionHeightmap(assets))
	r.HandleFunc("/{coatlicue}/{region}/watermap.r16", getRegionWatermap(assets))
}

func getLevelsListFunc(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := level.ListLevelIndex(assets)
		serveJson(result, w)
	}
}

func getLevelInfoFunc(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coatName := nwfs.NormalizePath(vars["coatlicue"])

		dir := level.NewCoatlicueDirectory(coatName)
		if !dir.Exists(assets) {
			http.NotFound(w, r)
			return
		}

		result := dir.LoadInfo(assets)
		if result == nil {
			http.NotFound(w, r)
			return
		}

		serveJson(result, w)
	}
}

func getRegionInfoFunc(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coatName := nwfs.NormalizePath(vars["coatlicue"])
		regionName := nwfs.NormalizePath(vars["region"])
		directory := level.NewRegionDirectory(coatName, regionName)

		info := directory.LoadRegionInfo(assets)
		if info == nil {
			http.NotFound(w, r)
			return
		}

		serveJson(info, w)
	}
}

func getRegionHeightmap(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coatName := nwfs.NormalizePath(vars["coatlicue"])
		regionName := nwfs.NormalizePath(vars["region"])
		directory := level.NewRegionDirectory(coatName, regionName)

		heightmap := directory.LoadHeightMap(assets)
		if heightmap == nil {
			http.NotFound(w, r)
			return
		}

		result := heightmap.ToFloat16Bytes()
		if result == nil {
			http.NotFound(w, r)
			return
		}

		serveContent(result, w, "application/octet-stream")
	}
}

func getRegionWatermap(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coatName := nwfs.NormalizePath(vars["coatlicue"])
		regionName := nwfs.NormalizePath(vars["region"])
		directory := level.NewRegionDirectory(coatName, regionName)

		watermap := directory.LoadWaterMap(assets)
		if watermap == nil {
			http.NotFound(w, r)
			return
		}

		result := watermap.ToFloat16Bytes()
		if result == nil {
			http.NotFound(w, r)
			return
		}

		serveContent(result, w, "application/octet-stream")
	}
}

func getRegionCapitalsFunc(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coatName := nwfs.NormalizePath(vars["coatlicue"])
		regionName := nwfs.NormalizePath(vars["region"])
		directory := level.NewRegionDirectory(coatName, regionName)

		result := directory.LoadRuntimeCapitals(assets)
		if result == nil {
			http.NotFound(w, r)
			return
		}

		serveJson(result, w)
	}
}
