import {
  Component,
  ElementRef,
  Input,
  OnInit,
  QueryList,
  ViewChild,
  ViewChildren,
} from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Conversation, Message, Order, User, UserType } from '@app/_models';
import {
  ConversationsService,
  OrdersService,
  UserService,
} from '@app/_services';
import { ChatService } from '@app/_services/chat.service';
import { HistoryService } from '@app/_services/history.service';
import { CrudDialogAction } from '@app/catalog/catalog.component';
import { OrderDialogComponent } from './order-dialog/order-dialog.component';

export interface OrderDialogData {
  action: CrudDialogAction;
  order?: Order;
  prevMessageId?: number;
  conversation: Conversation;
}

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss'],
})
export class ChatComponent implements OnInit {
  private _conversation!: Conversation;
  get conversation(): Conversation {
    return this._conversation;
  }

  public get CrudDialogAction() {
    return CrudDialogAction;
  }

  public get UserType() {
    return UserType;
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
  orders: Record<number, Order> = {};
  form = new FormGroup({
    message: new FormControl(''),
  });

  constructor(
    private historyService: HistoryService,
    private chatService: ChatService,
    userService: UserService,
    private conversationService: ConversationsService,
    public dialog: MatDialog,
    private orderService: OrdersService
  ) {
    this.user = userService.userValue!;
    orderService.orders.subscribe((orders) => (this.orders = orders));
    orderService.getAll().subscribe({
      error: (e) => console.log(e),
    });
  }

  ngOnInit(): void {
    this.historyService.messages.subscribe((messages) => {
      this.messages = messages[this.conversation.id];
      this.conversationService.setReadMessage(this.conversation.id);
    });
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

  openOrderDialog(data: OrderDialogData) {
    const dialogRef = this.dialog.open(OrderDialogComponent, {
      width: '500px',
      data,
    });
    dialogRef.afterClosed().subscribe((result: OrderDialogData) => {
      if (result && result.order) {
        this.orders[result.order.messageId] = result.order;
        this.chatService.sendOrder(result.order);
      }
    });
  }
}
