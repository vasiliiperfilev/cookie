import { Component } from '@angular/core';

import { Conversation, User } from '@app/_models';
import { UserService } from '@app/_services';

@Component({
  templateUrl: 'home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  user: User | null;
  currentConversation: Conversation | undefined;

  constructor(private userService: UserService) {
    this.user = this.userService.userValue;
  }

  selectConversation(c: Conversation) {
    this.currentConversation = c;
  }
}
