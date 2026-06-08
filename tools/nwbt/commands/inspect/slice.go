package inspect

import (
	"fmt"
	"io"
	"nw-buddy/tools/game"
	"nw-buddy/tools/nwfs"
	"nw-buddy/tools/rtti/nwt"
	"text/tabwriter"
)

type SliceInspector struct {
	count                  int
	MeshComponent          map[string]int
	SkinnedMeshComponent   map[string]int
	InstancedMeshComponent map[string]int
	AoiComponent           map[string]int
	FogVolumeComponent     map[string]int
	GameRigidBodyComponent map[string]int
	LightComponent         map[string]int
	MeshMergerComponent    map[string]int
	RiverComponent         map[string]int
	RoadComponent          map[string]int
	SnowComponent          map[string]int
	SplineComponent        map[string]int
	TimeOfDayPOIComponent  map[string]int
}

func NewSliceInspector() *SliceInspector {
	return &SliceInspector{
		MeshComponent:          make(map[string]int),
		SkinnedMeshComponent:   make(map[string]int),
		InstancedMeshComponent: make(map[string]int),
		AoiComponent:           make(map[string]int),
		FogVolumeComponent:     make(map[string]int),
		GameRigidBodyComponent: make(map[string]int),
		LightComponent:         make(map[string]int),
		MeshMergerComponent:    make(map[string]int),
		RiverComponent:         make(map[string]int),
		RoadComponent:          make(map[string]int),
		SnowComponent:          make(map[string]int),
		SplineComponent:        make(map[string]int),
		TimeOfDayPOIComponent:  make(map[string]int),
	}
}
func (it *SliceInspector) Inspect(assets *game.Assets, file nwfs.File) {
	it.count += 1

	f, err := assets.LoadSliceComponent(file)
	if err != nil {
		return
	}
	for _, entity := range f.Entities.Element {
		for _, component := range entity.Components.Element {

			switch v := component.(type) {
			case nwt.MeshComponent:
				options := v.Static_Mesh_Render_Node.Render_Options
				it.MeshComponent["Count"] += 1
				if v.Static_Mesh_Render_Node.Visible {
					it.MeshComponent["Visible"] += 1
				}
				if v.Load_Mesh_On_Activate {
					it.MeshComponent["Load_Mesh_On_Activate"] += 1
				}
				if options.VisibilityOccluder {
					it.MeshComponent["VisibilityOccluder"] += 1
				}
				if options.UseVisAreas {
					it.MeshComponent["UseVisAreas"] += 1
				}
				if options.ReceiveWind {
					it.MeshComponent["ReceiveWind"] += 1
				}
				if options.UseInProxyGeo {
					it.MeshComponent["UseInProxyGeo"] += 1
				}
				if options.UseManualViewDistance {
					it.MeshComponent["UseManualViewDistance"] += 1
				}

			case nwt.InstancedMeshComponent:
				options := v.Instanced_mesh_render_node.BaseClass1.Render_Options
				it.InstancedMeshComponent["Count"] += 1
				if options.VisibilityOccluder {
					it.InstancedMeshComponent["VisibilityOccluder"] += 1
				}
				if options.UseVisAreas {
					it.InstancedMeshComponent["UseVisAreas"] += 1
				}
				if options.ReceiveWind {
					it.InstancedMeshComponent["ReceiveWind"] += 1
				}
				if options.UseInProxyGeo {
					it.InstancedMeshComponent["UseInProxyGeo"] += 1
				}
				if options.UseManualViewDistance {
					it.InstancedMeshComponent["UseManualViewDistance"] += 1
				}

			case nwt.SkinnedMeshComponent:
				options := v.Skinned_Mesh_Render_Node.Render_Options
				it.SkinnedMeshComponent["Count"] += 1
				if v.Load_Mesh_On_Activate {
					it.SkinnedMeshComponent["Load_Mesh_On_Activate"] += 1
				}
				if options.UseVisAreas {
					it.SkinnedMeshComponent["UseVisAreas"] += 1
				}
				if options.NeverFrustumCull {
					it.SkinnedMeshComponent["NeverFrustumCull"] += 1
				}

			case nwt.AoiComponent:
				it.AoiComponent["Count"] += 1
				if v.M_useUserDefinedSpawnRadius {
					it.AoiComponent["M_useUserDefinedSpawnRadius"] += 1
				}
				if v.M_isStaticSlice {
					it.AoiComponent["M_isStaticSlice"] += 1
				}
				if v.M_overridewithuserdefinedspawnradius {
					it.AoiComponent["M_overridewithuserdefinedspawnradius"] += 1
				}
			case nwt.FogVolumeComponent:
				it.FogVolumeComponent["Count"] += 1
				if v.FogVolumeConfiguration.UseGlobalFogColor {
					it.FogVolumeComponent["UseGlobalFogColor"] += 1
				}
				if v.FogVolumeConfiguration.AffectsThisAreaOnly {
					it.FogVolumeComponent["AffectsThisAreaOnly"] += 1
				}
			case nwt.GameRigidBodyComponent:
				it.GameRigidBodyComponent["Count"] += 1
				if v.M_isDynamic {
					it.GameRigidBodyComponent["M_isDynamic"] += 1
				}
			case nwt.LightComponent:
				it.LightComponent["Count"] += 1
				config := v.LightConfiguration
				lightType := fmt.Sprintf("LightType%v", config.LightType)
				it.LightComponent[lightType] += 1
				if config.Ambient {
					it.LightComponent["Ambient"] += 1
				}
				if config.Deferred {
					it.LightComponent["Deferred"] += 1
				}
				if config.TerrainShadows {
					it.LightComponent["TerrainShadows"] += 1
				}
				if config.Visible {
					it.LightComponent["Visible"] += 1
				}
			case nwt.MeshMergerComponent:
				it.MeshMergerComponent["Count"] += 1
			case nwt.RiverComponent:
				it.RiverComponent["Count"] += 1
			case nwt.RoadComponent:
				it.RoadComponent["Count"] += 1
			case nwt.SnowComponent:
				it.SnowComponent["Count"] += 1
			case nwt.SplineComponent:
				it.SplineComponent["Count"] += 1
				config := v.Configuration
				splineType := fmt.Sprintf("SplineType%v", config.Spline_Type)
				it.SplineComponent[splineType] += 1
			case nwt.TimeOfDayPOIComponent:
				it.TimeOfDayPOIComponent["Count"] += 1
				if v.Configuration.HasTimeOfDayOverride {
					it.TimeOfDayPOIComponent["HasTimeOfDayOverride"] += 1
				}
			}
		}
	}

}

func (it *SliceInspector) Print(w io.Writer) {

	fmt.Fprintf(w, "Slices\t%d\n", it.count)

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintf(tw, "\nMeshComponent\n")
	for key, count := range it.MeshComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nInstancedMeshComponent\n")
	for key, count := range it.InstancedMeshComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nSkinnedMeshComponent\n")
	for key, count := range it.SkinnedMeshComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nAoiComponent\n")
	for key, count := range it.AoiComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nFogVolumeComponent\n")
	for key, count := range it.FogVolumeComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nGameRigidBodyComponent\n")
	for key, count := range it.GameRigidBodyComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nLightComponent\n")
	for key, count := range it.LightComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nMeshMergerComponent\n")
	for key, count := range it.MeshMergerComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nRiverComponent\n")
	for key, count := range it.RiverComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nRoadComponent\n")
	for key, count := range it.RoadComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nSnowComponent\n")
	for key, count := range it.SnowComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nSplineComponent\n")
	for key, count := range it.SplineComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}

	fmt.Fprintf(tw, "\nTimeOfDayPOIComponent\n")
	for key, count := range it.TimeOfDayPOIComponent {
		fmt.Fprintf(tw, "%s\t%d\n", key, count)
	}
	tw.Flush()

}
