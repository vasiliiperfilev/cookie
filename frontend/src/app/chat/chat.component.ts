import {
  Component,
  ElementRef,
  Input,
  OnInit,
  QueryList,
  ViewChild,
  ViewChildren,
} from '@angular/core';
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
  private _conversation!: Conversation; // private property _item

  // use getter setter to define the property
  get conversation(): Conversation {
    return this._conversation;
  }

  @Input({ required: true })
  set conversation(val: Conversation) {
    this._conversation = val;
    this.historyService.getMessagesByConversationId(val.id).subscribe({
      error: (err) => console.log(err),
    });
  }

  @ViewChild('chat', { read: ElementRef }) chatEl!: ElementRef;
  @ViewChildren('messages') messagesEl!: QueryList<any>;
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
      (messages) => (this.messages = messages[this.conversation.id])
    );
    this.historyService
      .getMessagesByConversationId(this.conversation.id)
      .subscribe({
        error: (err) => console.log(err),
      });
  }

  ngAfterViewInit() {
    this.scrollToBottom();
    this.messagesEl.changes.subscribe(this.scrollToBottom);
  }

  sendMessage() {
    if (!this.form.value.message) {
      return;
    }
    this.chatService.sendMessage(this.form.value.message, this.conversation.id);
    this.form.reset();
  }

  getSender(message: Message) {
    return this.conversation.users.find((user) => user.id === message.senderId)
      ?.name;
  }

  scrollToBottom = () => {
    try {
      this.chatEl.nativeElement.scrollTop =
        this.chatEl.nativeElement.scrollHeight;
    } catch (err) {
      console.log(err);
    }
  };
}
