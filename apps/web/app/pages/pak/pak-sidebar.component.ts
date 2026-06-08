import { Component, computed, inject } from '@angular/core'
import { LayoutModule } from '~/ui/layout'

import { rxResource, toSignal } from '@angular/core/rxjs-interop'
import { ActivatedRoute, Router } from '@angular/router'
import { debounceTime, map, of } from 'rxjs'
import { FileTreeComponent } from '~/ui/file-tree'
import { QuicksearchModule, QuicksearchService } from '~/ui/quicksearch'
import { FileTreeNode } from '../../ui/file-tree/file-tree.store'
import { PropertyGridModule } from '../../ui/property-grid'
import { PakService } from './pak.service'

@Component({
  standalone: true,
  selector: 'nwb-assets-sidebar',
  host: {
    class: 'ion-page',
  },
  imports: [LayoutModule, FileTreeComponent, QuicksearchModule, PropertyGridModule],
  providers: [
    QuicksearchService.provider({
      queryParam: 'search',
    }),
  ],
  template: `
    <ion-header class="bg-base-300">
      <ion-toolbar class="px-2">
        <nwb-quicksearch-input />
      </ion-toolbar>
    </ion-header>
    <div class="h-full flex flex-col">
      @if (files.isLoading()) {
        <div class="h-full flex items-center justify-center">
          <span class="loading loading-spinner"></span>
        </div>
      } @else {
        <nwb-file-tree
          class="px-2 flex-1"
          [files]="files.value()"
          [search]="search()"
          [selection]="selection()"
          (selected)="handleFileSelection($event)"
          (open)="handleOpen($event)"
        />
        <div class="flex-none h-8 font-bold bg-base-300 px-2 flex flex-row items-center">File Stat</div>
        <div class="flex-none h-64 overflow-auto p-2 bg-base-200">
          @if (fileStatResource.isLoading()) {
            <div class="h-full flex items-center justify-center">
              <span class="loading loading-spinner"></span>
            </div>
          } @else if (fileStatResource.error()) {
            <div class="text-error">
              {{ fileStatResource.error() }}
            </div>
          } @else if (fileStat()) {
            <nwb-property-grid [item]="fileStat()" class="text-xs" />
          }
        </div>
      }
    </div>
  `,
})
export class PakSidebarComponent {
  private service = inject(PakService)

  protected search = toSignal(inject(QuicksearchService).query$.pipe(debounceTime(300)))
  protected files = rxResource({
    stream: () => this.service.fetchFileList('**'),
  })

  private router = inject(Router)
  private route = inject(ActivatedRoute)

  protected selection = toSignal(this.route.queryParams.pipe(map((it) => it['file'])))
  protected fileStatResource = rxResource({
    params: this.selection,
    stream: ({ params }) => (params ? this.service.fetchFileStat(params) : of(null)),
  })
  protected fileStat = computed(() => flattenObject(this.fileStatResource.value()?.[0]))

  protected handleFileSelection(file: FileTreeNode) {
    if (file.isDir) {
      return
    }
    this.router.navigate(['.'], {
      queryParams: { file: file.id },
      queryParamsHandling: 'merge',
      relativeTo: this.route,
    })
  }

  protected handleOpen(file: FileTreeNode) {
    if (file.isDir) {
      return
    }
    const path = this.router
      .createUrlTree(['.'], {
        queryParams: { file: file.id },
        queryParamsHandling: 'merge',
        relativeTo: this.route,
      })
      .toString()
    window.open(path, '_blank')
  }
}
function flattenObject(obj: Record<string, any>, separator = '.'): Record<string, any> {
  const result: Record<string, any> = {}

  function recurse(current: any, path: string) {
    if (current === null || current === undefined) {
      result[path] = current
    } else if (typeof current !== 'object' || Array.isArray(current)) {
      result[path] = current
    } else {
      for (const key of Object.keys(current)) {
        recurse(current[key], path ? `${path}${separator}${key}` : key)
      }
    }
  }

  recurse(obj, '')
  return result
}
