import { Routes } from '@angular/router'

export const ROUTES: Routes = [
  {
    path: '',
    loadComponent: () => import('./encounter-page.component').then((it) => it.EncounterPageComponent),
    children: [
      {
        path: ':id',
        loadComponent: () => import('./encounter-detail-page.component').then((it) => it.EncounterDetailPageComponent),
      },
    ],
  },
]
