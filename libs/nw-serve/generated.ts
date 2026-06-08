/* Do not change, this code is generated from Golang structs */


export interface ViewerSlice {
  entities: ViewerEntity[];
  spawnRadius: number;
  isStaticSlice: boolean;
}
export interface RegionCapitalsData {
  capitals: {[key: string]: CapitalRuntimeData[]};
  chunks: {[key: string]: ChunkRuntimeData[]};
  slices: {[key: string]: ViewerSlice};
}
export interface RegionMaterialLayer {
  material: any;
  splatMap: any;
  affectedTiles: string;
  priority: number;
}
export interface RegionMaterial {
  tileX: number;
  tileY: number;
  defaultMaterial: any;
  normalMap: any;
  colorMap: any;
  specularMap: any;
  layers: RegionMaterialLayer[];
  pertinentLayersMipChain: string[];
}
export interface Vec2 {
  x: number;
  y: number;
}
export interface RegionImpostor {
  position: Vec2;
  model: string;
}
export interface RegionInfo {
  name: string;
  poiImpostors: RegionImpostor[];
  impostors: RegionImpostor[];
  terrainMaterial?: RegionMaterial;
}
export interface ViewerEntity {
  id: string;
  name: string;
  parentId?: string;
  transform: number[];
  components: any[];
}
export interface TimeOfDayVariable {
  name: string;
  color: string;
  value: string;
}
export interface TimeOfDay {
  time: number;
  timeStart: number;
  timeEnd: number;
  timeAnimSpeed: number;
  variables: TimeOfDayVariable[];
}
export interface RegionLocation {
  name: string;
  location: number[];
  playable: boolean;
}
export interface Location {
  x: number;
  y: number;
}
export interface ForcedRegion {
  location: Location;
  regionName: string;
}
export interface Region {
  name: string;
  spawnManifests: string[];
}
export interface Color {
  r: number;
  g: number;
  b: number;
}
export interface Tract {
  name: string;
  mapCategory: string;
  displayColor: Color;
}
export interface World {
  type: string;
  width: number;
  height: number;
}
export interface Document {
  tractmapCellSize: number;
  heightmapCellSize: number;
  regionSize: number;
  territoryMasterSlicePath: string;
  world: World;
  tracts: Tract[];
  regions: Region[];
  forcedRegions: ForcedRegion[];
}
export interface CoatlicueInfo {
  level: string;
  name: string;
  regionSize: number;
  regionCellSize: number;
  enableChunks: boolean;
  enableVegetation: boolean;
  enableDistribution: boolean;
  mountainHeight: number;
  mountainRoughness: number;
  oceanLevel: number;
  valleyIntensity: number;
  tracts?: Document;
  regions: RegionLocation[];
  gameModeMaps: GameModeMap[];
  missionTimeOfDay?: TimeOfDay;
  missionEntities: ViewerEntity[];
}
export interface GameModeMap {
  gameModeMapId: string;
  gameModeId: string;
  slicePath: string;
  coatlicueName: string;
  worldBounds: string[];
  teamTeleportData: string;
  uiMapId: string;
  sliceExclusionList: string[];
}
export interface CoatlicueListEntry {
  name: string;
  level: string;
  maps: GameModeMap[];
}
export interface LevelListEntry {
  name: string;
  coatlicueNames: string[];
}
export interface LevelIndex {
  levels: LevelListEntry[];
  coatlicues: CoatlicueListEntry[];
}
export interface AssetId {
  guid: string;
  subId: number;
}
export interface ServeAssetIdResult {
  asset?: Asset;
  assets: Asset[];
  link: AssetId;
  legacy: AssetId;
}
export interface Asset {
  guid: string;
  subId: number;
  type: string;
  file: string;
  size: number;
}
export interface ServeStatResultEntry {
  file?: string;
  asset?: Asset;
  dds?: {[key: string]: any};
}
export interface ServeStatResult {
  items: ServeStatResultEntry[];
}
export interface ServeListResult {
  items: string[];
}
export interface ServeApi {
  '/list/{filePattern}': ServeListResult;
  '/stats/{filePattern}': ServeStatResult;
  '/assets/{assetId}': ServeAssetIdResult;
  '/levels/list.json': LevelIndex;
  '/levels/{coatlicue}/info.json': CoatlicueInfo;
  '/levels/{coatlicue}/{region}/info.json': RegionInfo;
  '/levels/{coatlicue}/{region}/capitals.json': RegionCapitalsData;
  '/levels/{coatlicue}/{region}/heightmap.r16': ArrayBuffer;
  '/levels/{coatlicue}/{region}/watermap.r16': ArrayBuffer;
}







export interface AssetReference {
  guid: string;
  subId: number;
  hint?: string;
}
export interface CapitalRuntimeData {
  id: string;
  transform: number[];
  radius: number;
  slice: AssetReference;
}
export interface ChunkRuntimeData {
  id: string;
  transform: number[];
  size: number;
  slice: AssetReference;
}