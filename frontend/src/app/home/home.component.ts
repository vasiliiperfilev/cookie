import { Component } from '@angular/core';

import { User } from '@app/_models';
import { UserService } from '@app/_services';

@Component({
  templateUrl: 'home.component.html',
})
export class HomeComponent {
  user: User | null;
  currentConversationId: number | undefined;

  constructor(private userService: UserService) {
    this.user = this.userService.userValue;
  }

  selectConversation(id: number) {
    this.currentConversationId = id;
  }
}
