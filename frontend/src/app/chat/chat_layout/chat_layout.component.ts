import { Component } from '@angular/core';

import { Conversation, User } from '@app/_models';
import { ChatService, UserService } from '@app/_services';

@Component({
  templateUrl: 'chat_layout.component.html',
  styleUrls: ['chat_layout.component.scss'],
})
export class ChatLayoutComponent {
  user: User | null;
  currentConversation: Conversation | undefined;

  constructor(
    private userService: UserService,
    private chatService: ChatService
  ) {
    this.user = this.userService.userValue;
  }

  selectConversation(c: Conversation) {
    this.currentConversation = c;
  }
}
