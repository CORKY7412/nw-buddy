import { Component, computed, inject } from '@angular/core'
import { toSignal } from '@angular/core/rxjs-interop'
import { ActivatedRoute, RouterModule } from '@angular/router'
import { map } from 'rxjs'
import { LayoutModule } from '~/ui/layout'
import { SplitGutterComponent, SplitPaneDirective } from '~/ui/split-container'
import { PakDocsComponent } from './pak-docs.component'
import { PakNavbarComponent } from './pak-navbar.component'
import { PakPreviewCodeComponent } from './pak-preview-code.component'
import { PakPreviewImageComponent } from './pak-preview-image.component'
import { PakPreviewModelComponent } from './pak-preview-model.component'
import { PakPreviewSliceComponent } from './pak-preview-slice.component'
import { PakSidebarComponent } from './pak-sidebar.component'
import { PakService } from './pak.service'

@Component({
  standalone: true,
  selector: 'nwb-assets-page',
  providers: [PakService],
  imports: [
    PakSidebarComponent,
    LayoutModule,
    RouterModule,
    SplitPaneDirective,
    SplitGutterComponent,
    PakDocsComponent,
    PakNavbarComponent,

    PakPreviewCodeComponent,
    PakPreviewImageComponent,
    PakPreviewModelComponent,
    PakPreviewSliceComponent,
  ],
  host: {
    class: 'ion-page flex flex-row',
  },
  templateUrl: './pak-page.component.html',
})
export class PakPageComponent {
  protected service = inject(PakService)
  protected route = inject(ActivatedRoute)
  protected file = toSignal(this.route.queryParams.pipe(map((params) => params['file'] as string)))
  protected source = computed(() => this.service.fileSource(this.file()))
  protected viewerType = computed(() => this.source()?.viewer)
}
