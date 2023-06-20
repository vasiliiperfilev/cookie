import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { User } from '@app/_models';
import { Conversation, ConversationDto } from '@app/_models/conversation';
import { ConversationsService, UserService } from '@app/_services';

@Component({
  selector: 'app-conversations',
  templateUrl: './conversations.component.html',
  styleUrls: ['./conversations.component.scss'],
})
export class ConversationsComponent implements OnInit {
  @Output() selectConversationEvent = new EventEmitter<Conversation>();
  loading = false;
  conversations: Conversation[] = [];
  user: User;
  constructor(
    private conversationService: ConversationsService,
    userService: UserService
  ) {
    this.user = userService.userValue!;
  }

  ngOnInit() {
    this.loading = true;
    this.conversationService
      .getConversations()
      .subscribe({ error: (err) => console.log(err) });
    this.conversationService.conversations.subscribe(
      (conversations) => (this.conversations = conversations)
    );
  }

  getInitials(name: string) {
    const initials = name.charAt(0) + name.charAt(1);
    return initials.toUpperCase();
  }

  addConversation() {
    const c: ConversationDto = {
      userIds: [3, 4],
    };
    this.conversationService
      .postConversation(c)
      .subscribe({ error: (err) => console.log(err) });
  }

  selectConversation(c: Conversation) {
    this.selectConversationEvent.emit(c);
  }
}
