import {
  CoatlicueInfo,
  LevelIndex,
  RegionCapitalsData,
  RegionInfo,
  ServeAssetIdResult,
  ServeListResult,
  ServeStatResult,
} from './generated'

export type NwbtRequest<T> = {
  __type: T
  url: string
  type: 'json' | 'arraybuffer'
}

function nwbtJsonRequest<T>(url: string): NwbtRequest<T> {
  return {
    url,
    type: 'json',
  } satisfies Omit<NwbtRequest<T>, '__type'> as NwbtRequest<T>
}

function nwbtDataRequest(url: string): NwbtRequest<ArrayBuffer> {
  return {
    url,
    type: 'arraybuffer',
  } satisfies Omit<NwbtRequest<any>, '__type'> as NwbtRequest<ArrayBuffer>
}

export async function nwbtFetch<T>(baseUrl: string, req: NwbtRequest<T>): Promise<T> {
  if (baseUrl && baseUrl.endsWith('/')) {
    baseUrl = baseUrl.slice(0, -1)
  }
  if (!req.url.startsWith('/')) {
    req.url = '/' + req.url
  }
  return fetch((baseUrl || '') + req.url).then((it) => {
    if (req.type === 'json') {
      return it.json()
    }
    return it.arrayBuffer()
  })
}

export function nwbtFileUrl(filePath: string) {
  return `/files/${filePath}`
}

export function nwbtFileListUrl(pattern?: string) {
  if (pattern) {
    return nwbtJsonRequest<ServeListResult>(`/list/${pattern}`)
  }
  return nwbtJsonRequest<ServeListResult>(`/list`)
}

export function nwbtCatalogAssetUrl(assetId: string) {
  return nwbtJsonRequest<ServeAssetIdResult>(`/assets/${assetId}`)
}

export function nwbtFileStatUrl(filePattern: string) {
  return nwbtJsonRequest<ServeStatResult>(`/stats/${filePattern}`)
}

export function nwbtLevelsListUrl() {
  return nwbtJsonRequest<LevelIndex>(`/levels/list.json`)
}

export function nwbtLevelsCoatlicueInfoUrl(coatlicue: string) {
  return nwbtJsonRequest<CoatlicueInfo>(`/levels/${coatlicue}/info.json`)
}

export function nwbtLevelsCoatlicueRegionInfoUrl(coatlicue: string, region: string) {
  return nwbtJsonRequest<RegionInfo>(`/levels/${coatlicue}/${region}/info.json`)
}

export function nwbtLevelsCoatlicueRegionCapitalsUrl(coatlicue: string, region: string) {
  return nwbtJsonRequest<RegionCapitalsData>(`/levels/${coatlicue}/${region}/capitals.json`)
}

export function nwbtLevelsCoatlicueRegionHeightmapUrl(coatlicue: string, region: string) {
  return nwbtDataRequest(`/levels/${coatlicue}/${region}/heightmap.r16`)
}

export function nwbtLevelsCoatlicueRegionWatermapUrl(coatlicue: string, region: string) {
  return `/levels/${coatlicue}/${region}/watermap.r16`
}
