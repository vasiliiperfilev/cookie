import { Component, OnInit } from '@angular/core';
import { Conversation } from '@app/_models/conversation';
import { ConversationsService, UserService } from '@app/_services';

@Component({
  selector: 'app-conversations',
  templateUrl: './conversations.component.html',
  styleUrls: ['./conversations.component.sass'],
})
export class ConversationsComponent implements OnInit {
  loading = false;
  conversations: Conversation[] = [];
  constructor(private conversationService: ConversationsService) {}

  ngOnInit() {
    this.loading = true;
    this.conversationService.getConversations().subscribe((conversations) => {
      this.conversations = conversations;
      this.loading = false;
    });
  }

  getInitials(name: string) {
    const initials = name.charAt(0) + name.charAt(1);
    return initials.toUpperCase();
  }

  addConversation() {
    console.log('invoked');
    const c: Conversation = {
      id: 1,
      userIds: [1, 2],
    };
    this.conversationService
      .postConversation(c)
      .subscribe((conversation) => console.log(conversation));
  }
}
