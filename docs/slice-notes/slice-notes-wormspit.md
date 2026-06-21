```
16bs_sandworm.capitals.json
    "sliceName": "CoatGen\\\\513fe13\\worldevent_sandwormglimpse_spawner_master_2370482020.dynamicslice",

coatgen\513fe13\worldevent_sandwormglimpse_spawner_master_2370482020.dynamicslice.json
    AreaSpawnerComponent
        m_aliasasset
            "hint": "sharedassets/genericassets/aliasassets/randomencounter/16bs_brimstonesands/alias_randomencounter_sandworm.aliasasset"
```

Mention the next random encounter. It decides between two different versions of the sandworm.

```
sharedassets/genericassets/aliasassets/randomencounter/16bs_brimstonesands/alias_randomencounter_sandworm.aliasasset
    slices[0].percent: 100
    slices[0].slice.hint: "slices/worldevents/brimstonesands/sandwormbreach/enc_sandworm_breacharcevent_01.dynamicslice"
    slices[1].percent: 100
    slices[1].slice.hint: "slices/worldevents/brimstonesands/sandwormglimpse/enc_sandworm_glimpseevent_00.dynamicslice"
```

The actual ancounter, with a spawn timeline

```
slices/worldevents/brimstonesands/sandwormbreach/enc_sandworm_breacharcevent_01.dynamicslice
    EncounterComponent
        m_spawntimeline[]
            m_sliceasset
                "hint": "slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_launcher_burst.dynamicslice"
```

Sandworm can spit a projectile, we have the m_ammoid for the datasheet lookup

```
slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_launcher_burst.dynamicslice
    ProjectileSpawnerComponent
        "m_ammoid": "Sandworm_BreachEvent_AcidBallLauncher_Projectile_Burst"
```

Datasheet: javelindata_itemdefinitions_ammo.json

```
  {
   "AmmoID": "Sandworm_BreachEvent_AcidBallLauncher_Projectile_Burst",
   "AmmoType": "Siege",
   "DamageModifier": 1,
   "StaggerDamageModifier": 1,
   "AmmoPrefabPath": "WorldEvents/BrimstoneSands/SandwormBreach/Projectiles/Sandworm_BreachEvent_AcidBall_Projectile_Burst",
   "AmmoWhizByTrigger": "Play_Bullet_WizzBys"
  }
```

brings us to the projectile slice

```
slices/WorldEvents/BrimstoneSands/SandwormBreach/Projectiles/Sandworm_BreachEvent_AcidBall_Projectile_Burst.dynamicslice
    ProjectileComponent
        m_serverfacetptr
            m_spawnonhitasset
                "hint": "slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_projectile_burst_aoe.dynamicslice"
```

on hit, it bursts and turns into an Ore (indirection via PrefabSpawner)

```
slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_projectile_burst_aoe.dynamicslice
    PrefabSpawnerComponent
        m_sliceasset
            "hint": "slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_chunk_00.dynamicslice"

slices/worldevents/brimstonesands/sandwormbreach/projectiles/sandworm_breachevent_acidball_chunk_00.dynamicslice
    PrefabSpawnerComponent
        m_aliasasset
            "hint": "sharedassets/genericassets/aliasassets/randomencounter/16bs_brimstonesands/sandwormbreachevent/alias_sandwormbreachevent_ore_lg.aliasasset"
```

the kind of ore is chosen by the last aliasasset. all percentages are at 100, but we still have different chances for different ores, based on how many entries of each ore there are. in this case we have

- 5 entries of sandstone
- 2 entries of brimstone
- 1 entry for gypsum ore

```
sharedassets/genericassets/aliasassets/randomencounter/16bs_brimstonesands/sandwormbreachevent/alias_sandwormbreachevent_ore_lg.aliasasset.json
    slices[0].percent: 100
    slices[0].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_sandstone_single_a_spawn.dynamicslice"
    slices[1].percent: 100
    slices[1].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_sandstone_single_a_spawn.dynamicslice"
    slices[2].percent: 100
    slices[2].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_sandstone_single_a_spawn.dynamicslice"
    slices[3].percent: 100
    slices[3].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_sandstone_single_a_spawn.dynamicslice"
    slices[4].percent: 100
    slices[4].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_sandstone_single_a_spawn.dynamicslice"
    slices[5].percent: 100
    slices[5].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_brimstone_single_a_spawn.dynamicslice"
    slices[6].percent: 100
    slices[6].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_brimstone_single_a_spawn.dynamicslice"
    slices[7].percent: 100
    slices[7].slice.hint: "slices/worldevents/brimstonesands/randomencounter/io/io_randomencounter_single_gypsumore_a_spawns.dynamicslice"

```
