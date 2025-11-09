import { CommonModule } from '@angular/common'
import { ChangeDetectionStrategy, Component, inject } from '@angular/core'
import { ActivatedRoute } from '@angular/router'
import { NwModule } from '~/nw'
import { LayoutModule } from '~/ui/layout'
import { observeRouteParam } from '~/utils'
import { EncounterDetailComponent } from '../../../widgets/data/encounter-detail'
import { toSignal } from '@angular/core/rxjs-interop'

@Component({
  selector: 'nwb-encounter-detail-page',
  templateUrl: './encounter-detail-page.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, NwModule, LayoutModule, EncounterDetailComponent],
  host: {
    class: 'ion-page',
  },
})
export class EncounterDetailPageComponent {
  protected route = inject(ActivatedRoute)
  protected recordId = toSignal(observeRouteParam(this.route, 'id'))
}
