import { Component, inject } from '@angular/core'
import { PakService } from './pak.service'

@Component({
  selector: 'nwb-pak-docs',
  templateUrl: './pak-docs.component.html',
  host: {
    class: 'ion-page items-center',
  },
})
export class PakDocsComponent {
  protected service = inject(PakService)
}
