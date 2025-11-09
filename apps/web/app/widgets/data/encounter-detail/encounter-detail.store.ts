import { computed, inject, linkedSignal } from '@angular/core'
import { patchState, signalMethod, signalStore, withComputed, withMethods, withState } from '@ngrx/signals'
import { describeNodeSize } from '@nw-data/common'
import { Feature, FeatureCollection, MultiPoint } from 'geojson'
import { injectNwData } from '~/data'
import { resourceValue } from '~/utils'
import { GameMapService } from '../../game-map'

export interface EncounterDetailState {
  encounterId: string
}
export type EncounterFeature = Feature<MultiPoint, EncounterFeatureProperties>
export type EncounterFeatureCollection = FeatureCollection<MultiPoint, EncounterFeatureProperties>
export interface EncounterFeatureProperties {
  id: string
  color: string
  label: string
  title: string
}

export const EncounterDetailStore = signalStore(
  withState<EncounterDetailState>({
    encounterId: null,
  }),
  withMethods((state) => {
    return {
      connect: signalMethod((encounterId: string) => {
        patchState(state, {
          encounterId,
        })
      }),
    }
  }),
  withComputed(({ encounterId }) => {
    const db = injectNwData()
    const record = resourceValue({
      defaultValue: null,
      keepPrevious: true,
      params: encounterId,
      loader: ({ params }) => {
        return db.encounterMeta(params)
      },
    })
    return {
      record,
    }
  }),
  withComputed(({ record }) => {
    const mapService = inject(GameMapService)
    return {
      name: computed(() => record()?.name),
      tag: computed(() => record()?.tag),
      stages: computed(() => record()?.stages || []),
      mapIds: computed(() => {
        const spawns = record()?.spawns || []
        const ids: string[] = []
        for (const spawn of spawns) {
          if (!ids.includes(spawn.mapID)) {
            ids.push(spawn.mapID)
          }
        }
        return ids.sort()
      }),
      mapFeatures: computed(() => {
        const result: Record<string, EncounterFeatureCollection> = {}
        const props = describeNodeSize('Medium')
        let featureId = 0
        const spawns = record()?.spawns || []
        for (const { mapID, positions } of spawns || []) {
          result[mapID] ||= {
            type: 'FeatureCollection',
            features: [],
          }
          result[mapID].features.push({
            id: featureId++,
            type: 'Feature',
            geometry: {
              type: 'MultiPoint',
              coordinates: positions.map((it) => mapService.xyToLngLat(it)),
            },
            properties: {
              id: null,
              color: props.color,
              label: null,
              title: null,
            },
          })
        }
        return result
      }),
    }
  }),
  withComputed(({ mapIds, mapFeatures }) => {
    const mapId = linkedSignal(() => mapIds()?.[0])
    const mapData = computed(() => mapFeatures()?.[mapId()])
    const mapBounds = computed(() => selectBounds(mapData()))
    return {
      mapId,
      mapData,
      mapBounds,
    }
  }),
)

function selectBounds(data: EncounterFeatureCollection): [number, number, number, number] {
  if (!data) {
    return null
  }
  let min: [number, number] = null
  let max: [number, number] = null
  for (const feature of data.features) {
    for (const [x, y] of feature.geometry.coordinates) {
      if (!min) {
        min = [x, y]
        max = [x, y]
      } else {
        min[0] = Math.min(min[0], x)
        min[1] = Math.min(min[1], y)
        max[0] = Math.max(max[0], x)
        max[1] = Math.max(max[1], y)
      }
    }
  }
  if (min[0] === max[0] || min[1] === max[1]) {
    min[0] -= 0.001
    min[1] -= 0.001
    max[0] += 0.001
    max[1] += 0.001
  }
  return [min[0], min[1], max[0], max[1]]
}
