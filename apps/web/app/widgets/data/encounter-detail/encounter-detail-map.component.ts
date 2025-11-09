import { Component, inject, input, viewChild } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { NwModule } from '~/nw'
import { IconsModule } from '~/ui/icons'
import { svgExpand } from '~/ui/icons/svg'
import { TooltipModule } from '~/ui/tooltip'
import { GameMapComponent, GameMapCoordsComponent, GameMapLayerDirective } from '~/widgets/game-map'
import { EncounterDetailStore } from './encounter-detail.store'

@Component({
  selector: 'nwb-encounter-detail-map',
  templateUrl: './encounter-detail-map.component.html',
  imports: [
    NwModule,
    TooltipModule,
    FormsModule,
    IconsModule,
    GameMapComponent,
    GameMapLayerDirective,
    GameMapCoordsComponent,
  ],
  host: {
    class: 'block relative',
    '[class.hidden]': '!isVisible',
  },
})
export class EncounterItemDetailMapComponent {
  protected store = inject(EncounterDetailStore)
  protected iconExpand = svgExpand

  protected mapComponent = viewChild(GameMapComponent)

  protected get isVisible() {
    return !!this.store.mapId()
  }
}
