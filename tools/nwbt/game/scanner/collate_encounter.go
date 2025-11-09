package scanner

import (
	"nw-buddy/tools/utils/maps"
	"strings"
)

type ScannedEncounter struct {
	EncounterID string                  `json:"encounterID"`
	Name        string                  `json:"name"`
	Tag         string                  `json:"tag"`
	Stages      []EncounterStage        `json:"stages"`
	Spawns      []ScannedEncounterSpawn `json:"spawns"`
}

type ScannedEncounterSpawn struct {
	MapID     string     `json:"mapID"`
	Positions []Position `json:"positions"`
}

func CollateEncounter(rows []EncounterEntry) (result []ScannedEncounter, count int) {
	result = make([]ScannedEncounter, 0)
	index := maps.NewDict[*maps.Dict[*ScannedEncounterSpawn]]()
	records := maps.NewDict[ScannedEncounter]()
	for _, row := range rows {
		records.Store(row.EncounterID, ScannedEncounter{
			EncounterID: row.EncounterID,
			Name:        row.Name,
			Tag:         row.Tag,
			Stages:      row.Stages,
		})

		mapId := strings.ToLower(row.MapID)
		position := PositionFromV3(row.Position).Truncate()
		recordID := strings.ToLower(row.EncounterID)

		node := index.
			LoadOrCreate(recordID, maps.NewDict).
			LoadOrCreate(mapId, func() *ScannedEncounterSpawn {
				return &ScannedEncounterSpawn{
					MapID: mapId,
				}
			})
		node.Positions = append(node.Positions, position)
	}

	for recordID, b1 := range index.SortedIter() {
		record := records.Get(recordID)
		record.Spawns = make([]ScannedEncounterSpawn, 0)
		for _, value := range b1.SortedIter() {

			positions := sortAndFilterPositions(value.Positions)
			count += len(positions)
			record.Spawns = append(record.Spawns, ScannedEncounterSpawn{
				MapID:     value.MapID,
				Positions: positions,
			})
		}
		result = append(result, record)
	}

	return
}
