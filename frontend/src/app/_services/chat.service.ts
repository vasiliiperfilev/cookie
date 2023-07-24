import { Injectable } from '@angular/core';
import {
  Message,
  Order,
  WsEventType,
  WsMessageEvent,
  WsOrderEvent,
} from '@app/_models';
import { environment } from '@environments/environment';
import { Subject } from 'rxjs';
import { webSocket } from 'rxjs/webSocket';
import { HistoryService } from './history.service';
import { OrdersService } from './order.service';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  private wsConn: Subject<WsMessageEvent | WsOrderEvent>;

  constructor(
    private historyService: HistoryService,
    private userService: UserService,
    private orderService: OrdersService
  ) {
    this.wsConn = webSocket<WsMessageEvent | WsOrderEvent>(
      `${environment.webSocketUrl}/v1/chat?token=${userService.tokenValue?.token}`
    );
    this.wsConn.subscribe({
      next: (e) => this.receiveEvent(e), // Called whenever there is a message from the server.
      error: (err) => console.log(err), // Called if at any point WebSocket API signals some kind of error.
      complete: () => console.log('complete'), // Called when connection is closed (for whatever reason).
    });
  }

  sendMessage(content: string, conversationId: number) {
    const msgs = this.historyService.messagesValue[conversationId];
    const prevMessageId =
      msgs && msgs.length > 0 ? msgs[msgs.length - 1].id : 0;
    const wsMsgEvt: WsMessageEvent = {
      type: WsEventType.MESSAGE,
      payload: {
        conversationId: conversationId,
        prevMessageId: prevMessageId,
        content: content,
      },
    };
    this.wsConn.next(wsMsgEvt);
  }

  sendOrder(order: Order) {
    const wsMsgEvt: WsOrderEvent = {
      type: WsEventType.NEW_ORDER,
      payload: order,
    };
    this.wsConn.next(wsMsgEvt);
  }

  sendUpdatedOrder(order: Order) {
    const wsMsgEvt: WsOrderEvent = {
      type: WsEventType.UPDATE_ORDER,
      payload: order,
    };
    this.wsConn.next(wsMsgEvt);
  }

  private receiveEvent(evt: WsMessageEvent | WsOrderEvent) {
    if (evt.type === WsEventType.MESSAGE) {
      this.historyService.pushToLocalHistory(evt.payload as Message);
    } else if (
      evt.type === WsEventType.NEW_ORDER ||
      evt.type === WsEventType.UPDATE_ORDER
    ) {
      this.historyService
        .getMessagesById((evt.payload as Order).messageId)
        .subscribe({
          error: (err) => console.log(err),
        });
      this.orderService.pushToLocalOrders(evt.payload as Order);
    }
  }
}
