import { Component, inject } from '@angular/core'
import { DomSanitizer } from '@angular/platform-browser'
import { env } from '../../../environments/env'
@Component({
  selector: 'nwb-level-viewer',
  template: `<iframe class="w-full h-full" [src]="iframeUrl"></iframe>`,
  host: {
    class: 'ion-page',
  },
})
export class ViewerComponent {
  private sanitizer = inject(DomSanitizer)
  protected iframeUrl = this.sanitizer.bypassSecurityTrustResourceUrl(env.nwbtUrl)
}
