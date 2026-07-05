import { Component, computed, inject, input } from '@angular/core'
import { DomSanitizer } from '@angular/platform-browser'
import { FileSource, PakService } from './pak.service'
import { env } from '../../../environments/env'

@Component({
  selector: 'nwb-pak-preview-model',
  template: ` <iframe class="w-full h-full bg-base-300" [attr.src]="iframeUrl()"></iframe> `,
  host: {
    class: 'ion-page',
  },
})
export class PakPreviewModelComponent {
  private sanitizer = inject(DomSanitizer)
  public service = inject(PakService)
  public source = input<FileSource>()

  public iframeUrl = computed(() => {
    const url = new URL(env.nwbtUrl)
    url.searchParams.set('model', this.source().file)
    return this.sanitizer.bypassSecurityTrustResourceUrl(url.toString())
  })
}
