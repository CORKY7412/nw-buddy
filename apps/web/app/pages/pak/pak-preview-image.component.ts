import { Component, computed, inject, input } from '@angular/core'
import { FileSource, PakService } from './pak.service'

@Component({
  selector: 'nwb-pak-preview-image',
  template: ` <img [src]="imageUrl()" /> `,
})
export class PakPreviewImageComponent {
  public service = inject(PakService)
  public source = input<FileSource>()
  public imageUrl = computed(() => this.source().base + this.source().file + '.png')
}
