package serve

import (
	"net/http"
	"nw-buddy/tools/formats/dds"
	"nw-buddy/tools/game"
	"nw-buddy/tools/nwfs"
)

func GetListHandler(archive nwfs.Archive) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pattern := nwfs.NormalizePath(r.URL.Path)
		list, err := archive.Glob(pattern)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		files := make([]string, len(list))
		for i, file := range list {
			files[i] = file.Path()
		}
		serveJson(files, w)
	}
}

func GetStatHandler(assets *game.Assets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pattern := nwfs.NormalizePath(r.URL.Path)
		list, err := assets.Archive.Glob(pattern)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		result := make([]map[string]any, len(list))
		for i, file := range list {
			filePath := file.Path()

			stat := make(map[string]any)
			if asset := assets.Catalog.FindByFile(file.Path()); asset != nil {
				stat["asset"] = asset
			} else {
				stat["file"] = filePath
			}

			if dds.IsDDS(filePath) || dds.IsDDSAlpha(filePath) {
				meta, _ := dds.LoadMeta(file)
				if meta != nil {
					stat["dds"] = meta.Stats()
				}
			}

			result[i] = stat
		}

		serveJson(result, w)
	}
}
