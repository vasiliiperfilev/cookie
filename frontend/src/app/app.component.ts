import { Component } from '@angular/core';

import { User, UserType } from './_models';
import { UserService } from './_services';

@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  styleUrls: ['app.component.scss'],
})
export class AppComponent {
  user?: User | null;

  constructor(private userService: UserService) {
    this.userService.user.subscribe((x) => (this.user = x));
  }

  public get UserType() {
    return UserType;
  }

  logout() {
    this.userService.logout();
  }
}
