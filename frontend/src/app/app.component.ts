import { Component } from "@angular/core";

import { UserService } from "./_services";
import { User } from "./_models";

@Component({ selector: "app-root", templateUrl: "app.component.html" })
export class AppComponent {
  user?: User | null;

  constructor(private userService: UserService) {
    this.userService.user.subscribe((x) => (this.user = x));
  }

  logout() {
    this.userService.logout();
  }
}
