# NW Level Format

this file contains notes about the level file structure in new world and how nw-buddy currently parses it.

## File structure

```
[ROOT]
├─ ...
├─ levels                          // natural place for cry engine levels
│  │                               // mostly abandoned for new world, except for main menu
│  └─[LVL_NAME]
│     ├─ levelinfo.xml             // references mission file e.g mission0
│     ├─ mission_mission0.xml      // environment (lighting, fog, sky) setup. probably just defaults
│     ├─ mission0.entities_xml     // in case of main menu this contains a full scene. otherwise mostly empty
│     ├─ resourcelist.txt          // compilation artifact, lists all resources in this level
│
├─ sharedassets/coatlicue/[LVL_NAME]
│  ├─ regions
│  │  └─r_+[XX]_+[YY]                 // r_+00_+00 always present, open world has more
│  │    ├─ capitals
│  │    │  └─ [LAYER]                  // artist defined layer name
│  │    │     └─ [LAYER].capitals.json // listing of dynamic slice spawns
│  │    │     └─ [LAYER].metadata      // either a build artifact or actual metadata about this layer
│  │    ├─ impostors                  // impostor cell models directory
│  │    │  ├─ impostor_0.cgf          // 256 impostor cells per region, 16x16 grid,
│  │    │  ├─ impostor_0.cgfheap      // each cell is 128m x 128m covering 2048m x 2.048m region
│  │    │  ├─ ...
│  │    │  ├─ impostor_255.cgf
│  │    │  └─ impostor_255.cgfheap
│  │    ├─ poi_impostors              // impostor cell models directory
│  │    │  ├─ impostor_0.cgf          // 256 impostor cells per region, 16x16 grid,
│  │    │  ├─ impostor_0.cgfheap      // each cell is 128m x 128m covering 2048m x 2.048m region
│  │    │  ├─ ...
│  │    │  ├─ impostor_255.cgf
│  │    │  └─ impostor_255.cgfheap
│  │    ├─ localmappings.json  // contains tract colors as seen in tractmap.tif file
│  │    ├─ mapsettings.json    // contains region size (2048)
│  │    ├─ impostors.json      // low level geometry index list, trees and bushes
│  │    ├─ poi_impostors.json  // low level geometry index list, rocks and structures
│  │    ├─ region.chunks       // same purpose as capitals, quad tree optimized
│  │    ├─ region.dstribution  // gatherables and mass object distribution
│  │    ├─ region.heightmap    // 16-bit single channel tiff
│  │    ├─ region.tractmap.tif // rgba color coded regions (mountains, roads, rivers, etc)
│  │    ├─ region.vegetation   // not parsed yet, expected vegetation coverage data
│  │    └─ region.waterqt      // water quadtree, essentially a heightmap for water surfaces
│  ├─ offlineoptions.json      // contains the impostor cell size (128)
│  ├─ playable.json            // list of playable regions
│  ├─ terrain.json             // references terrain material (splat maps), ocean level, mountain height
│  └─ tracts.json              // world size, region size, tracts, masterSlicePath
```

## `tracts.json`

Most levels have the very same content of this file. So unfortunately it contains nothing unique to a level.

nw-buddy doesn't really use this

```json
{
  "tractmapCellSize": 8,
  "heightmapCellSize": 1,
  "regionSize": 2048,
  // all levels have this very same entry, except open world.
  "territoryMasterSlicePath": "slices/POIs/Territories/Territories_Master_Combat.dynamicslice",

  "world": {
    "type": "Island",
    "width": 1, // number of regions in x direction
    "height": 1 // number of regions in y direction
  },

  "tracts": [
    { "name": "blank", "displayColor": { "r": 0, "g": 153, "b": 0 } },
    { "name": "starterarea", "displayColor": { "r": 102, "g": 51, "b": 102 } },
    { "name": "grassland", "displayColor": { "r": 204, "g": 255, "b": 204 } }
  ],

  "regions": [{ "name": "greenZone", "spawnManifests": ["greenZone"] }],

  "forcedRegions": [{ "location": { "x": 0, "y": 0 }, "regionName": "greenZone" }]
}
```

## `impostors.json` and `poi_impostors.json`

Both are index list files that reference the actual 3D geometry in the `impostors/` and `poi_impostors/` folders

```json
{
  "materialAssetID": "",
  "impostors": [
    {
      "cellIndex": 0,
      // directly references the .cgf asset
      "meshAssetID": "{D46366CA-BD6A-5CBE-8722-E3D7E067ED37}",
      // absolute world position, model is placed at that specific location without adjustments.
      // vertical offset is baked into the model, so z is always 0.0
      "worldPosition": {
        "x": 64.0,
        "y": 64.0
      }
    }
    // ...
  ]
}
```

## `offlineoptions.json`

```json
{
  "distributionBakeSeed": 93016, // unsure for what the seed is
  "distributionDataEnabled": true, // probably enables usage of region.distribution file
  "vegetationDataEnabled": true, // probably enables usage of region.vegetation file

  "streaming": {
    "impostorCellEdgeLength": 128,
    "regionChunksEnabled": true // probably enables usage of region.chunks files
  }
}
```

## `world.json` and `playable.json`

`world.json` lists all regions in this level. `playable.json` lists only playable regions, which are a subset of world.json.
Probably used to allocate server resources.

```json
[
  [0, 0] // region address [X, Y], e.g. [0, 5] -> is outpost rush map on top left
]
```

## `terrain.json`

```json
{
  "generatorType": "Heightmap",
  "heightCrop": 0.0,
  "mountainRoughness": 0.699999988079071,
  // mountainHeight of 256 can be treated as "no terrain"
  // levels with actual terrain always have 2048, or in case of ftue it's 1024
  "mountainHeight": 2048.0,
  "snowMinimumSlope": 0.0,
  "snowStartHeight": 100.0,
  "valleyIntensity": 0.1694214940071106,
  // ocean level value of -1000 means no ocean
  "oceanLevel": 40.0,
  // terrain rendering material, splat maps
  "worldMaterialAssetPath": "Materials/terrain/NW_OPR_004_Trench/NW_OPR_004_Trench.worldmat"
}
```

## `region.distribution`

Mass placement of gatherables and other objects

Binary data

json converted example:

```json
{
  // Region Index
  "region": [0, 0],
  // list of slices. prefix and file extension are missing
  "slices": [
    "",
    "gatherables/master_bush", // -> slices/gatherables/master_bush.{dynamicslice,slice.meta}
    "gatherables/master_bush"
    // ...
  ],
  // The slices define the base gatherable type.
  // this list defines the actual variant to spawn with the slice at same index.
  "variants": [
    "",
    "Bush_WaxMyrtle_a",
    "Bush_WaxMyrtle_b"
    // ...
  ],
  // glue index between positions and slices/variants, see code below for details
  "indices": [
    723, 711, 712
    // ...
  ],
  //
  "positions": [
    [49967, 60559],
    [49743, 60543],
    [50055, 60743]
    // ...
  ],
  // not sure what these are, we ignore them
  "positions2": [],
  "types2": "",
  "positions3": [],
  "types3": ""
}
```

format parser at: https://github.com/giniedp/nw-buddy/blob/live/tools/nwbt/formats/distribution/format.go

usage example: https://github.com/giniedp/nw-buddy/blob/live/tools/nwbt/game/scanner/scan_distributions.go

Positions need to be decoded

- into region space coordinates (0-2048)
- and then into world space coordinates, by adding the region offset (region index \* 2048)

```go
// converts the position value from region.destribution into world space coordinates
// - region: region address [X, Y] (e.g. [0, 5])
// - position: position to decode
func ConvertPosition(region [2]uint32, position [2]uint16) (float32, float32, float32, float32) {
	areaSize := float32(2048) // should be a parameter, but it's always 2048 for our regions
	maxValue := float32(0xFFFF)
	rx := float32(region[0]) * areaSize
	ry := float32(region[1]) * areaSize

	px := (float32(position[0]) / maxValue) * areaSize
	py := (float32(position[1]) / maxValue) * areaSize
	x := rx + px
	y := ry + py
	return x, y, px, py
}
```

## `*.capitals.json`

Entry point for dynamic slice placement.
Each entry has a position/rotation where a specific slice should be spawned.

```json
{
  "capitals": [
    {
      "id": "6f218469-dae9-618c-c7ec-67caf94962eb",
      // absolute world position, no adjustments.
      // all contents of the loaded slice are relative to this position/rotation
      "worldPosition": {
        "x": 220.0,
        "y": 930.0,
        "z": 143.0
      },
      "rotation": {
        "x": 0.0,
        "y": 0.0,
        "z": 0.0,
        "w": 1.0
      },
      "footprint": {
        "type": "Circle",
        "id": "3ce007f9-9711-0778-dadd-1f878cffbb20",
        "radius": 1.0 // always 1.0, too small to be meaningful
      },
      "sliceName": "CoatGen\\831a57a1\\AI_CutNav_Box_3m_6f218469dae9618cc7ec67caf94962eb",
      "sliceAssetId": "{E4520C92-6BE0-536E-8E1F-2050EDB6CC91}:4962eb"
    }
    // ...
  ]
}
```

When slices are registered, search for the `AoiComponent` in the root entity of the slice. This holds the necessary parameters for streaming the slice in and out.

```json
{
  "__type": "AoiComponent",
  "baseclass1": {
    "__type": "FacetedComponent",
    "baseclass1": {
      "__type": "AZ::Component",
      "id": "12556981330868048132"
    },
    "m_replicationindex": 0
  },
  "m_aoigridcategory": 6,
  "m_aoiradius": 0,
  "m_slicephysicalgridradius": 0.20000000298023224,
  "m_slicedetectiongridradius": 0,
  "m_slicephysicalradius": 0,
  "m_additionalslicephysicalminradius": 0,
  "m_isstaticslice": true,
  "m_slicetags": 0,
  "m_slicespawnradius": 2083.581298828125, // most meaningful parameter
  "m_useuserdefinedspawnradius": false,
  "m_overridewithuserdefinedspawnradius": false,
  "m_editorsliceviewradius": 0,
  "m_editorslicephysicalradius": 0,
  "m_editorslicespawnradius": 0,
  "m_editoraoiradius": 0,
  "m_editorisstaticslice": false,
  "m_editorwillimpostor ": false,
  "m_editorisrequiredonserver": false,
  "m_editorrefreshbutton": false
},
```

## `region.chunks` and `[CAPITAL_NAME].chunks`

Chunks are only present in open world level. The `*.capitals.json` still exist, but the referenced slice name and asset id are striped like this:

```
  "sliceName": "<PLOT>",
  "sliceAssetId": "{00000000-0000-0000-0000-000000000000}:0"
```

For each capital, there then is a corresponding `[CAPITAL_NAME].chunks` file, which contains a quad-tree optimized list of chunk placements.
Functionality it's the same as capitals.json, but with a quad tree structure and certainly more optimized streaming the large amount of open world data.

```json
{
  "__type": "AC608BE6-77F3-5AF5-A7A9-607621389D91",
  "chunks": {
    "__type": "283F62D6-A310-5D1E-A38E-409DB6C165A4",
    "element": [
      {
        "__type": "ChunkEntry",
        "cellindex": {
          "__type": "CellIndex",
          "x": "0",
          "y": "0",
          "z": "0"
        },
        "size": "1024",
        "spawnradius": 0,
        "layer": "08qp_roads",
        "chunktype": 0,
        // the absolute world position where the slice should spawn
        "worldposition": [6656, 6656, 48.77204513549805],
        // the asset to spawn
        "assetid": {
          "__type": "AssetId",
          "guid": "ebaf58f8-cd54-5f1e-a90f-ae4e3ddda43d",
          "subid": 3197738318
        }
      }
    ]
  }
}
```

the `region.chunks` at region level seem to combine all the `*.capitals.chunks` entries in one file. This is unverified though, but some samples I checked seem to confirm this.

## Traversing level files

Historically, nw-buddy does it a bit fuzzy and flaky, since it doesn't replicate the actual streaming behavior and is more interested in finding all the "possible" spawns in the level, even if they are gated behind certain conditions.

Buddy globs the following files

```
  "**/region.distribution",                     // 91
  "**/coatlicue/**/regions/**/*.capitals.json", // 1362
  "**.dynamicslice",                            // 196350
  "!lyshineui/**",
```

For each it detects

- the file type (capitals, distribution, slice)
- the level it's in

and then walks the encountered capitals/slices recursively. Results are grouped by level

The old "per file" handling is here

- https://github.com/giniedp/nw-buddy/blob/live/tools/nwbt/game/scanner/scan.go

The recent refactor added the following package, that tries to be more aligned with the level structure

- https://github.com/giniedp/nw-buddy/tree/live/tools/nwbt/game/level

Walking the slice a simple tree traversal, but there are 2 important hierarchies to be aware of:

1. The Entity hierarchy inside the slice. We can ignore that for the entity position, since all entities come with their absolute world position already resolved. I haven't found any slice where this is not the case. This absolute position is always relative to where the slice was spawned.
2. The slcie spawn hierarchy. When a slice spawns at a certain position, that position is the reference point for all entities in that slice.

An example of a GameTransformComponent with hierarchy reference and absolute world position:

```json
{
  "__type": "GameTransformComponent",
  "baseclass1": {
    "__type": "FacetedComponent",
    "baseclass1": {
      "__type": "AZ::Component",
      "id": "10034558780467598761"
    },
    "m_replicationindex": 0
  },
  // the absolute world transform of this entity. But absolute in the sense of "in this slice"
  "m_worldtm": {

    "data": [
      -0.9781476259231567, 0.20791161060333252, 0,
      -0.20791161060333252, -0.9781476259231567, 0,
      0, 0, 1,
      -2.4482421875, 3.535888671875, -1.7242584228515625
    ]
  },
  "m_parentid": {
    "__type": "EntityId",
    "id": "16624699878889673125" // the parent entity ID. usually inside the slice
  },
  // we can ignore this local transform for parsing
  // at game runtime, when this changes, it should issue a world transform update
  // but in the files, nothing ever changes
  "m_localtm": {
    "data": [
      -0.9781476259231567, 0.20791161060333252, 0,
      -0.20791161060333252, -0.9781476259231567, 0,
      0, 0, 1,
      -2.4482421875, 3.535888671875, -1.7242584228515625
    ]
  },
  "m_onnewparentkeepworldtm": true,
  "m_isstatic": false
},
```

Besides the `GameTransformComponent` a hierarchy may also be formed by a `TransformComponent`
Functionality is the same, but different class different data layout

Resolving the transform for both cases is here

- https://github.com/giniedp/nw-buddy/blob/e76a0e8fce2f1bd4dfc3b669104714cd31364dac/tools/nwbt/game/utils.go#L125
- https://github.com/giniedp/nw-buddy/blob/e76a0e8fce2f1bd4dfc3b669104714cd31364dac/tools/nwbt/utils/math/mat4/mat4.go#L141

## Traversing slice spawns

nw-buddy walks deep into the entity/slice hierarchy by following each and every "spawner" it comes across

- https://github.com/giniedp/nw-buddy/blob/e76a0e8fce2f1bd4dfc3b669104714cd31364dac/tools/nwbt/game/scanner/scan_slice.go#L445

The spawners all have different characteristics. nw-Buddy tries to handle them all somehow.
While this works well for the website, there may be other systems involved at runtime that interact with the spawners differently.

- SpawnerComponent - just a base class, should't exists as standalone, but is still checked
- PointSpawnerComponent - spawns an asset at the position of the entity
- PrefabSpawnerComponent - same as point spawner. Most of static level geometry is from this spawner.
- ProjectileSpawnerComponent - spawns a projectile asset.
  - Needs ammo datasheet lookup to find the asset
  - https://www.nw-buddy.de/datasheets?file=javelindata_itemdefinitions_ammo.json
- ProjectileComponent - contains "spawnOnHitAsset", e.g. gleamite meteor turn into gleamite ore
- AreaSpawnerComponent - list of locations with spawn slices
- EncounterManagerComponent - contains phases of an encounter, and minion spawns during the phases

Here is some old note when i tried to trace the spawn of the sandworm. It goes through most of the spawner types.

- [./slice-notes/slice-notes-wormspit.md](./slice-notes/slice-notes-wormspit.md)

## Harvesting entity data

During the traversal, nw-buddy tries to gather all possible information about this "slice branch". Same source file as above

- https://github.com/giniedp/nw-buddy/blob/e76a0e8fce2f1bd4dfc3b669104714cd31364dac/tools/nwbt/game/scanner/scan_slice.go#L29

Depends on the slice combination, data may be provided at different levels of the hierarchy.
For most data, if provided at higher level, it overrides what is provided at lower (leaf) levels.
This logic is enough for the website scanner, but should have specific rules at runtime.

This is what gathered from specific components

- NpcComponent
  - NpcID
- VitalsComponent
  - VitalsID -> vitals datasheet entry
- ActionListComponent
  - DamageTable -> datasheet entry
  - AdbFile -> Animation Database file \*.adb, used for animation player
- AIVariantProviderComponent
  - VitalsID -> vitals datasheet entry
  - CategoryID -> vitals categories datasheet entry
  - Level -> creature spawn level
  - Info whether creature level should taken from territory level
- VariationDataComponent
  - VariantID: a variant for a gatherable. refers to a datasheet entry with visual and loot variations
- GatherableControllerComponent
  - GatherableID -> gatherable datasheet entry
- ReadingInteractionComponent
  - LoreID -> lore note datasheet entry
- SkinnedMeshComponent/MeshComponent
  - Mesh asset for 3d viewer
- HousingPlotComponent
  - HouseType: house tier types. no datasheet exists for this
- AssemblyComponent
  - StationID: crafting station type. no datasheet exists for this
