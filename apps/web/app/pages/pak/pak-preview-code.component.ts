import { httpResource } from '@angular/common/http'
import { Component, computed, inject, input } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { CodeEditorModule } from '../../ui/code-editor'
import { MonacoSliceExtensionDirective } from './monaco'
import { FileSource, PakService, textTypeMap, toJsonTypes, toLuaTypes } from './pak.service'

@Component({
  selector: 'nwb-pak-preview-code',
  template: `
    <nwb-code-editor
      class="ion-page"
      [nwbMonacoSliceExtension]=""
      (nwbMonacoAssetCommand)="service.handleMonacoAssetCommand($event)"
      [ngModel]="textContent.value()"
      [language]="fileInfo().lang"
      [disabled]="true"
    />
  `,
  host: {
    class: 'ion-page',
  },
  imports: [CodeEditorModule, MonacoSliceExtensionDirective, FormsModule],
})
export class PakPreviewCodeComponent {
  public service = inject(PakService)
  public source = input<FileSource>()

  protected fileInfo = computed(() => {
    const ext = this.source()?.ext

    if (toJsonTypes.includes(ext)) {
      return {
        url: this.source().base + this.source().file + '.json',
        lang: 'json',
      }
    }
    if (toLuaTypes.includes(ext)) {
      return {
        url: this.source().base + this.source().file + '.lua',
        lang: 'lua',
      }
    }
    return {
      url: this.source().base + this.source().file,
      lang: textTypeMap[ext] || 'txt',
    }
  })

  protected textContent = httpResource.text(() => this.fileInfo().url)
}
