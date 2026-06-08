import { HttpClient } from '@angular/common/http'
import { inject, Injectable } from '@angular/core'
import { nwbtFileListUrl, nwbtFileStatUrl, ServeListResult, ServeStatResult } from '@nw-serve'
import { environment } from 'apps/web/environments'
import { map } from 'rxjs'

const toImageTypes = ['dds', 'png', 'tif', 'a', '1a', '2a', '3a', '4a', '5a', '6a', '7a', 'heightmap']
const toModelTypes = ['cgf', 'cdf', 'skin', 'mtl', 'dynamicslice']
const toLuaTypes = ['luac']
const toJsonTypes = [
  'aliasasset',
  'chunks',
  'datasheet',
  'distribution',
  'dynamicslice',
  'dynamicuicanvas',
  'meta',
  'metadata',
  'refreshzoneconfigs',
  'slicedata',
  'timeline',
  'waterqt',

  'bsdb',
  'craftstationdb',
  'crestdb',
  'evnotdb',
  'fishdb',
  'fueldb',
  'gactdb',
  'paperdolldb',
  'pbadb',
  'radb',
  'rankdb',
  'sprd',
  'uidb',
  'aoffdb',
  'equipdb',
  'gmevtdb',
  'gdb',
  'gadb',
  'gds',
  'eyecolordb',
  'facemarkdb',
  'hairstyledb',
  'skintonedb',

  'collisionfilters',
]
const textTypeMap = {
  json: 'json',
  xml: 'xml',
  txt: 'txt',
  cfg: 'txt',
  ext: 'txt',

  mtl: 'xml',
  cdf: 'xml',
  chrparams: 'xml',
  animevents: 'xml',
  bspace: 'xml',
  comb: 'xml',
  adb: 'xml',
  grid: 'xml',
  actionlist: 'xml',
  entities_xml: 'xml',
  worldmat: 'xml',
  regionmat: 'xml',
  surfacemap: 'xml',
}

export interface FileSource {
  baseUrl: string

  path: string
  ext: string
  textPath?: string
  textType?: string
  modelPath?: string
  imagePath?: string
}

@Injectable({ providedIn: 'root' })
export class PakService {
  private http = inject(HttpClient)

  public fileSource(file: string) {
    if (!file) {
      return null
    }
    const basename = file.split('/').pop()
    const tokens = basename.split('.')
    let ext = tokens.pop()
    if (ext.match(/^[0-9]+$/)) {
      ext = tokens.pop()
    }
    const stat: FileSource = {
      baseUrl: this.nwbtUrl('file/'),
      path: file,
      ext: ext,
    }
    if (textTypeMap[ext]) {
      stat.textPath = file
      stat.textType = textTypeMap[ext]
    }
    if (toJsonTypes.includes(ext)) {
      stat.textPath = `${file}.json`
      stat.textType = 'json'
    }
    if (toLuaTypes.includes(ext)) {
      stat.textPath = `${file}.lua`
      stat.textType = 'lua'
    }
    if (toModelTypes.includes(ext)) {
      stat.modelPath = `${file}.glb?cache=0`
    }
    if (toImageTypes.includes(ext)) {
      stat.imagePath = `${file}.png`
    }
    return stat
  }

  public nwbtUrl(resource: string) {
    if (!resource.startsWith('/')) {
      resource = '/' + resource
    }
    return `${environment.nwbtUrl}${resource}`
  }

  public fileUrl(file: string, format?: string) {
    if (format) {
      file = `${file}.${format}`
    }
    return this.nwbtUrl(`file/${file}`)
  }

  public fileListUrl(pattern: string) {
    return this.nwbtUrl(nwbtFileListUrl(pattern).url)
  }

  public fileStatUrl(file: string) {
    return this.nwbtUrl(nwbtFileStatUrl(file).url)
  }

  public fetchFileList(pattern: string) {
    return this.http.get<ServeListResult>(this.fileListUrl(pattern)).pipe(map((result) => result.items))
  }

  public fetchFileStat(pattern: string) {
    return this.http.get<ServeStatResult>(this.fileStatUrl(pattern)).pipe(map((result) => result.items))
  }
}
