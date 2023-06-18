import { Injectable } from '@angular/core';
import { Message, WsEventType, WsMessageEvent } from '@app/_models';
import { Subject } from 'rxjs';
import { environment } from '@environments/environment';
import { webSocket } from 'rxjs/webSocket';
import { HistoryService } from './history.service';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  private wsConn: Subject<WsMessageEvent>;

  constructor(
    private historyService: HistoryService,
    userService: UserService
  ) {
    this.wsConn = webSocket<WsMessageEvent>(
      `${environment.webSocketUrl}/v1/chat?token=${userService.tokenValue?.token}`
    );
    this.wsConn.subscribe({
      next: (e) => this.receiveEvent(e), // Called whenever there is a message from the server.
      error: (err) => console.log(err), // Called if at any point WebSocket API signals some kind of error.
      complete: () => console.log('complete'), // Called when connection is closed (for whatever reason).
    });
  }

  sendMessage(content: string, conversationId: number) {
    const prevMessageId = this.historyService.messagesValue[
      this.historyService.messagesValue.length - 1
    ]
      ? this.historyService.messagesValue[
          this.historyService.messagesValue.length - 1
        ].id
      : 0;
    const wsMsgEvt: WsMessageEvent = {
      type: WsEventType.MESSAGE,
      payload: {
        conversationId: conversationId,
        prevMessageId: prevMessageId,
        content: content,
      },
    };
    this.wsConn.next(wsMsgEvt);
    this.historyService.pushToLocalHistory({
      ...wsMsgEvt.payload,
      createdAt: new Date(),
    } as Message);
  }

  private receiveEvent(evt: WsMessageEvent) {
    if (evt.type === WsEventType.MESSAGE) {
      this.historyService.pushToLocalHistory(evt.payload as Message);
    }
  }
}
