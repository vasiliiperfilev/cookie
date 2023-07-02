import { inject } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivateFn,
  Router,
  RouterStateSnapshot,
} from '@angular/router';
import { UserType } from '@app/_models';

import { UserService } from '@app/_services';

export const canActivate: CanActivateFn = (
  next: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const userService = inject(UserService);
  const router = inject(Router);
  const token = userService.tokenValue;
  if (token) {
    // authorised so return true
    return true;
  }

  // not logged in so redirect to login page with the return url
  router.navigate(['/login'], {
    queryParams: { returnUrl: state.url },
  });
  return false;
};

export const canActivateSupplier: CanActivateFn = (
  next: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const userService = inject(UserService);
  const router = inject(Router);
  const token = userService.tokenValue;
  if (token) {
    const user = userService.userValue;
    if (user?.type === UserType.SUPPLIER) {
      return true;
    }
  }

  return false;
};
