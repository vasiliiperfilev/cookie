import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { canActivate, canActivateSupplier } from './_helpers';
import { CatalogComponent } from './catalog/catalog.component';
import { ChatLayoutComponent } from './chat';

const accountModule = () =>
  import('./account/account.module').then((x) => x.AccountModule);

const routes: Routes = [
  { path: 'chat', component: ChatLayoutComponent, canActivate: [canActivate] },
  {
    path: 'catalog',
    component: CatalogComponent,
    canActivate: [canActivateSupplier],
  },
  { path: '', loadChildren: accountModule },

  // otherwise redirect to login
  { path: '**', redirectTo: 'login' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
