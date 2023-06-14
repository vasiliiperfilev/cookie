import { Component, Input, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Message, WsEventType, WsMessageEvent } from '@app/_models';
import { ChatService } from '@app/_services/chat.service';
import { HistoryService } from '@app/_services/history.service';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.sass'],
})
export class ChatComponent implements OnInit {
  @Input({ required: true }) conversationId!: number;
  messages: Message[] = [];
  form = new FormGroup({
    message: new FormControl('', [Validators.required]),
  });

  constructor(
    private historyService: HistoryService,
    private chatService: ChatService
  ) {}

  ngOnInit(): void {
    this.historyService.messages.subscribe(
      (messages) => (this.messages = messages)
    );
    this.historyService
      .getMessagesByConversationId(this.conversationId)
      .subscribe({
        error: (err) => console.log(err),
      });
  }

  sendMessage() {
    if (!this.form.value.message) {
      return;
    }
    this.chatService.sendMessage(this.form.value.message, this.conversationId);
    this.form.value.message = '';
  }
}
