import { ChangeDetectionStrategy, Component, computed, inject, input } from '@angular/core'
import { EncounterItemDetailMapComponent } from './encounter-detail-map.component'
import { EncounterDetailStore } from './encounter-detail.store'
import { JsonPipe } from '@angular/common'
import { CodeEditorComponent } from '../../../ui/code-editor'
import { FormsModule } from '@angular/forms'

@Component({
  selector: 'nwb-encounter-detail',
  template: `
    <nwb-encounter-detail-map />
    @if (data()) {
      <nwb-code-editor [ngModel]="data() | json" [language]="'json'" [disabled]="true" class="h-full" />
    }
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [EncounterItemDetailMapComponent, JsonPipe, CodeEditorComponent, FormsModule],
  providers: [EncounterDetailStore],
})
export class EncounterDetailComponent {
  public store = inject(EncounterDetailStore)

  public encounterId = input<string>(null)
  public data = computed(() => {
    if (!this.store.record()) {
      return null
    }
    return {
      name: this.store.name(),
      tag: this.store.tag(),
      stages: this.store.stages(),
    }
  })
  public constructor() {
    this.store.connect(this.encounterId)
  }
}
