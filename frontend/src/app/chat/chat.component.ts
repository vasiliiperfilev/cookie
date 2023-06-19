import { Component, Input, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Conversation, Message, User } from '@app/_models';
import { UserService } from '@app/_services';
import { ChatService } from '@app/_services/chat.service';
import { HistoryService } from '@app/_services/history.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss'],
})
export class ChatComponent implements OnInit {
  @Input({ required: true }) conversation!: Conversation;
  messages: Message[] = [];
  user: User;
  form = new FormGroup({
    message: new FormControl(''),
  });

  constructor(
    private historyService: HistoryService,
    private chatService: ChatService,
    private userService: UserService
  ) {
    this.user = userService.userValue!;
  }

  ngOnInit(): void {
    this.historyService.messages.subscribe(
      (messages) => (this.messages = messages)
    );
    this.historyService
      .getMessagesByConversationId(this.conversation.id)
      .subscribe({
        error: (err) => console.log(err),
      });
  }

  sendMessage() {
    if (!this.form.value.message) {
      return;
    }
    this.chatService.sendMessage(this.form.value.message, this.conversation.id);
    this.form.reset();
  }
}
