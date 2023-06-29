import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { UserService } from '@app/_services';

@Component({ templateUrl: 'layout.component.html' })
export class LayoutComponent {
  constructor(private router: Router, private userService: UserService) {
    // redirect to chat if already logged in
    if (this.userService.tokenValue) {
      this.router.navigate(['/chat']);
    } else {
      this.router.navigate(['/login']);
    }
  }
}
