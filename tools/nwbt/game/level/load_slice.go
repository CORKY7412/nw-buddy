package level

import (
	"fmt"
	"log/slog"
	"nw-buddy/tools/formats/catalog"
	"nw-buddy/tools/game"
	"nw-buddy/tools/nwfs"
	"nw-buddy/tools/rtti/nwt"
	"nw-buddy/tools/utils"
	"nw-buddy/tools/utils/math/mat4"
)

type TrackSliceAsset func(asset catalog.AssetId)

func LoadViewerSlice(assets *game.Assets, file nwfs.File, track TrackSliceAsset) (ViewerSlice, error) {
	if track == nil {
		track = func(asset catalog.AssetId) {}
	}
	result := ViewerSlice{
		Entities: make([]ViewerEntity, 0),
	}

	slice, err := assets.LoadSliceComponent(file)
	if err != nil {
		return result, err
	}

	if slice == nil {
		return result, nil
	}

	for _, entity := range slice.Entities.Element {
		viewerEntity := ViewerEntity{
			ID:         fmt.Sprintf("%d", entity.Id.Id),
			Name:       string(entity.Name),
			Transfrom:  game.FindEntityTransformMat4(&entity),
			Components: make([]any, 0),
		}
		if parent := game.FindEntityParent(slice, &entity); parent != nil {
			viewerEntity.ParentId = fmt.Sprintf("%d", parent.Id.Id)
		}

		for _, comp := range entity.Components.Element {
			switch v := comp.(type) {
			case nwt.AoiComponent:
				result.IsStaticSlice = bool(v.M_isStaticSlice)
				result.SpawnRadius = float32(v.M_sliceSpawnRadius)
			case nwt.MeshComponent:
				meshNode := v.Static_Mesh_Render_Node
				if game.ShouldIgnoreMeshNode(meshNode) {
					continue
				}
				modelAsset := meshNode.Static_Mesh
				modelFile, err := assets.LookupFileByAsset(modelAsset)
				if err != nil {
					slog.Warn("skip missing model", "err", err)
					continue
				}
				if modelFile == nil {
					continue
				}
				model := modelFile.Path()

				materialAsset := meshNode.Material_Override_Asset
				materialFile, err := assets.LookupFileByAsset(materialAsset)
				if err != nil {
					slog.Warn("ignore model material", "err", err)
				}
				material := ""
				if materialFile != nil {
					material = materialFile.Path()
				}

				model, material = assets.ResolveCgfAndMtl(model, material)
				viewerEntity.Components = append(viewerEntity.Components, ViewerMeshComponent{
					Type:                   MeshComponentName,
					Model:                  model,
					Material:               material,
					LoadOnActivate:         bool(v.Load_Mesh_On_Activate),
					Opacity:                float32(meshNode.Render_Options.Opacity),
					MaxViewDistance:        float32(meshNode.Render_Options.MaxViewDistance),
					ViewDistanceMultiplier: float32(meshNode.Render_Options.ViewDistanceMultiplier),
					CrossFadeTime:          float32(meshNode.Render_Options.CrossFadeTime),
					CastShadow:             bool(meshNode.Render_Options.CastShadows),
					ShouldInstance:         bool(meshNode.Render_Options.ShouldInstance),
					ShouldMerge:            bool(meshNode.Render_Options.ShouldMerge),
					ForceMerge:             bool(meshNode.Render_Options.ForceMerge),
					FadeEnabled:            bool(meshNode.Render_Options.FadeEnabled),
					AcceptDecals:           bool(meshNode.Render_Options.AcceptDecals),
					AcceptSand:             bool(meshNode.Render_Options.AcceptSand),
					AcceptSilhouette:       bool(meshNode.Render_Options.AcceptSilhouette),
					AcceptSnow:             bool(meshNode.Render_Options.AcceptSnow),
					AlwaysRender:           bool(meshNode.Render_Options.AlwaysRender),
					SortType:               int8(meshNode.Render_Options.SortType),
					VisibilityOccluder:     bool(meshNode.Render_Options.VisibilityOccluder),
					UseVisAreas:            bool(meshNode.Render_Options.UseVisAreas),
				})
			case nwt.InstancedMeshComponent:
				meshNode := v.Instanced_mesh_render_node.BaseClass1
				if game.ShouldIgnoreMeshNode(meshNode) {
					continue
				}
				modelAsset := meshNode.Static_Mesh
				modelFile, err := assets.LookupFileByAsset(modelAsset)
				if err != nil {
					slog.Warn("skip missing model", "err", err)
					continue
				}
				if modelFile == nil {
					continue
				}
				model := modelFile.Path()

				materialAsset := meshNode.Material_Override_Asset
				materialFile, err := assets.LookupFileByAsset(materialAsset)
				if err != nil {
					slog.Warn("ignore model material", "err", err)
				}
				material := ""
				if materialFile != nil {
					material = materialFile.Path()
				}

				instances := make([]mat4.Data, 0)
				for _, t := range v.Instanced_mesh_render_node.Instance_transforms.Element {
					instances = append(instances, mat4.FromAzTransform(t))
				}

				model, material = assets.ResolveCgfAndMtl(model, material)
				viewerEntity.Components = append(viewerEntity.Components, ViewerMeshComponent{
					Type:                   MeshComponentName,
					Model:                  model,
					Material:               material,
					Instances:              instances,
					LoadOnActivate:         true,
					Opacity:                float32(meshNode.Render_Options.Opacity),
					MaxViewDistance:        float32(meshNode.Render_Options.MaxViewDistance),
					ViewDistanceMultiplier: float32(meshNode.Render_Options.ViewDistanceMultiplier),
					CrossFadeTime:          float32(meshNode.Render_Options.CrossFadeTime),
					CastShadow:             bool(meshNode.Render_Options.CastShadows),
					ShouldInstance:         bool(meshNode.Render_Options.ShouldInstance),
					ShouldMerge:            bool(meshNode.Render_Options.ShouldMerge),
					ForceMerge:             bool(meshNode.Render_Options.ForceMerge),
					FadeEnabled:            bool(meshNode.Render_Options.FadeEnabled),
					AcceptDecals:           bool(meshNode.Render_Options.AcceptDecals),
					AcceptSand:             bool(meshNode.Render_Options.AcceptSand),
					AcceptSilhouette:       bool(meshNode.Render_Options.AcceptSilhouette),
					AcceptSnow:             bool(meshNode.Render_Options.AcceptSnow),
					AlwaysRender:           bool(meshNode.Render_Options.AlwaysRender),
					SortType:               int8(meshNode.Render_Options.SortType),
					VisibilityOccluder:     bool(meshNode.Render_Options.VisibilityOccluder),
					UseVisAreas:            bool(meshNode.Render_Options.UseVisAreas),
				})

			case nwt.SpawnerComponent:
				facet, ok := v.BaseClass1.M_serverFacetPtr.(nwt.SpawnerComponentServerFacet)
				if !ok {
					continue
				}
				list := make([]nwt.AzAsset, 0)
				list = utils.AppendUniqNoZero(list, facet.M_sliceAsset)
				list = utils.AppendUniqNoZero(list, facet.M_aliasAsset)
				asset := findAsset(assets, list)
				if asset == nil {
					continue
				}

				track(asset.AssetId)
				viewerEntity.Components = append(viewerEntity.Components, ViewerSpawnerComponent{
					Type: SpawnerComponentName,
					Slice: AssetReference{
						Guid:  asset.Guid,
						SubId: asset.SubID,
						Hint:  asset.File,
					},
				})
			case nwt.PointSpawnerComponent:
				list := make([]nwt.AzAsset, 0)
				list = utils.AppendUniqNoZero(list, v.BaseClass1.M_sliceAsset)
				list = utils.AppendUniqNoZero(list, v.BaseClass1.M_aliasAsset)
				asset := findAsset(assets, list)
				if asset == nil {
					continue
				}
				track(asset.AssetId)
				viewerEntity.Components = append(viewerEntity.Components, ViewerPointSpawnerComponent{
					Type: PointSpawnerComponentName,
					Slice: AssetReference{
						Guid:  asset.Guid,
						SubId: asset.SubID,
						Hint:  asset.File,
					},
				})
			case nwt.PrefabSpawnerComponent:
				list := make([]nwt.AzAsset, 0)
				list = utils.AppendUniqNoZero(list, v.M_sliceAsset)
				list = utils.AppendUniqNoZero(list, v.M_aliasAsset)
				asset := findAsset(assets, list)
				if asset == nil {
					continue
				}
				track(asset.AssetId)
				viewerEntity.Components = append(viewerEntity.Components, ViewerPrefabSpawnerComponent{
					Type: PrefabSpawnerComponentName,
					Slice: AssetReference{
						Guid:  asset.Guid,
						SubId: asset.SubID,
						Hint:  asset.File,
					},
				})
			case nwt.AreaSpawnerComponent:
				facet, ok := v.BaseClass1.M_serverFacetPtr.(nwt.AreaSpawnerComponentServerFacet)
				if !ok {
					break
				}
				list := make([]nwt.AzAsset, 0)
				list = utils.AppendUniqNoZero(list, facet.M_sliceAsset)
				list = utils.AppendUniqNoZero(list, facet.M_aliasAsset)
				asset := findAsset(assets, list)
				if asset == nil {
					continue
				}

				locations := make([]mat4.Data, 0)
				for _, loc := range facet.M_locations.Element {
					entity := game.FindEntityById(slice, loc.EntityId.Id)
					if entity == nil {
						continue
					}
					locations = append(locations, game.FindEntityTransformMat4(entity))
				}
				track(asset.AssetId)
				viewerEntity.Components = append(viewerEntity.Components, ViewerAreaSpawnerComponent{
					Type: AreaSpawnerComponentName,
					Slice: AssetReference{
						Guid:  asset.Guid,
						SubId: asset.SubID,
						Hint:  asset.File,
					},
					Locations:        locations,
					MinResspawnRange: float32(facet.M_minrespawnrange),
					MaxRespawnRange:  float32(facet.M_maxrespawnrange),
					LiveCount:        int(facet.M_livecount),
					SpawnOnEnable:    bool(facet.M_spawnonenable),
					SpawnOnTrigger:   bool(facet.M_spawnOnTrigger),
				})
			case nwt.SplineComponent:
				// TODO:
			case nwt.FogVolumeComponent:
				// TODO:
			case nwt.DecalComponent:
				// TODO:
			case nwt.GameRigidBodyComponent:
				// TODO:
			case nwt.HitVolumeComponent:
				// TODO:
			case nwt.LightComponent:
				// TODO:
			case nwt.MeshMergerComponent:
				// TODO:
			case nwt.RiverComponent:
				// TODO:
			case nwt.RoadComponent:
				// TODO:
			case nwt.SnowComponent:
				// TODO:
			case nwt.VegetationAreaComponent:
				// TODO:
			}
		}

		if len(viewerEntity.Components) > 0 {
			result.Entities = append(result.Entities, viewerEntity)
		}
	}

	return result, nil
}

func findAsset(assets *game.Assets, list []nwt.AzAsset) *catalog.Asset {
	if len(list) == 0 {
		return nil
	}

	isEmpty := true
	for _, asset := range list {
		if catalog.UUID(asset.Guid).IsZeroOrEmpty() {
			continue
		}

		isEmpty = false

		result := assets.Catalog.Find(asset.Guid, asset.Type, asset.Hint)
		if result != nil {
			return result
		}
	}

	if !isEmpty {
		slog.Warn("asset not found in catalog", "tried", list)
	}

	return nil
}
