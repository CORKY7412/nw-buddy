import { Component, computed, effect, ElementRef, inject, signal, untracked, viewChildren } from '@angular/core'
import { IconsModule } from '../../ui/icons'
import { svgCode, svgCubes, svgImage, svgXmark } from '../../ui/icons/svg'
import { PakService } from './pak.service'

@Component({
  selector: 'nwb-pak-navbar',
  template: `
    @if (tabs().length) {
      @for (tab of tabs(); track tab.id) {
        <a
          #tabButton
          role="tab"
          class="tab flex gap-1"
          [class.italic]="!tab.materialized"
          [class.tab-active]="tab === activeTab()"
          (click)="active.set(tab.id)"
          (dblclick)="tab.materialized = true"
        >
          <nwb-icon [icon]="tab.icon" class="w-4 h-4" />
          {{ tab.name }}
          @if (canClose()) {
            <nwb-icon [icon]="iconClose" (click)="handleClose(tab, $event)" class="w-4 h-4" />
          }
        </a>
      }
    }
  `,
  imports: [IconsModule],
  host: {
    role: 'tablist',
    class: 'tabs tabs-sm tabs-lift',
  },
})
export class PakNavbarComponent {
  private service = inject(PakService)
  protected element = inject(ElementRef).nativeElement
  protected buttons = viewChildren('tabButton', { read: ElementRef })
  protected iconClose = svgXmark
  protected tabs = signal<TabbarItem[]>([])
  protected active = signal<string>(null)
  protected activeTab = computed(() => {
    const activeId = this.active()
    return this.tabs().find((it) => it.id === activeId)
  })
  protected canClose = computed(() => this.tabs().length > 1)

  public constructor() {
    this.resolveState()
    effect(() => {
      const fileId = this.service.selectedFileId()
      untracked(() => this.open(fileId))
    })
    effect(() => {
      const active = this.activeTab()
      untracked(() => {
        if (!active || this.service.selectedFileId() === active.id) {
          return
        }
        this.service.selectFileId(active.id)
      })
    })
    effect(() => {
      this.tabs()
      this.active()
      untracked(() => {
        this.persistState()
      })
    })
  }

  public open(file: string) {
    if (!file) {
      this.active.set(null)
      return
    }

    const found = this.tabs().find((it) => it.id === file)
    if (found) {
      this.active.set(found.id)
      return
    }

    if (!this.activeTab() || this.activeTab().materialized) {
      this.tabs.update((tabs) => [...tabs, fileToTab(file)])
      this.active.set(file)
      return
    }

    this.tabs.update((tabs) => {
      return tabs.map((it) => {
        if (it.id === this.activeTab().id) {
          return fileToTab(file)
        }
        return it
      })
    })
    this.active.set(file)
  }

  protected handleClose(tab: TabbarItem, event: Event) {
    event.stopPropagation()
    const index = this.tabs().findIndex((it) => it.id === tab.id)
    this.tabs.update((tabs) => tabs.filter((it) => it.id !== tab.id))
    if (this.active() !== tab.id) {
      return
    }
    const nextTab = this.tabs()[index] || this.tabs()[index - 1]
    this.active.set(nextTab?.id || null)
  }

  protected resolveState() {
    try {
      const data = JSON.parse(localStorage.getItem('nw-buddy-pak-tabs')) as { files: string[]; active: string }
      if (data) {
        this.tabs.set(data.files.map(fileToTab).map((it) => ({ ...it, materialized: true })))
        if (!this.service.selectedFileId()) {
          this.active.set(data.active)
        }
      }
    } catch (e) {
      console.error('Failed to resolve state', e)
    }
  }

  protected persistState() {
    try {
      const data = {
        files: this.tabs()
          .filter((it) => it.materialized)
          .map((it) => it.id),
        active: this.active(),
      }
      localStorage.setItem('nw-buddy-pak-tabs', JSON.stringify(data))
    } catch (e) {
      console.error('Failed to persist state', e)
    }
  }
}

export interface TabbarItem {
  id: string
  name: string
  icon: string
  materialized: boolean
}

function fileToTab(file: string): TabbarItem {
  return {
    id: file,
    name: file.split('/').pop(),
    icon: fileToIcon(file),
    materialized: false,
  }
}

function fileToIcon(file: string): string {
  const ext = file.split('.').pop()
  switch (ext) {
    case 'dynamicslice':
      return svgCubes
    case 'mesh':
      return svgCubes
    case 'texture':
      return svgImage
    default:
      return svgCode
  }
}
