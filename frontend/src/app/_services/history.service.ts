import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Message } from '@app/_models/message';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { ConversationsService } from './conversations.service';

@Injectable({
  providedIn: 'root',
})
export class HistoryService {
  private messagesSubject: BehaviorSubject<Record<number, Message[]>>;
  public messages: Observable<Record<number, Message[]>>;

  constructor(
    private http: HttpClient,
    private conversationService: ConversationsService
  ) {
    this.messagesSubject = new BehaviorSubject<Record<number, Message[]>>({});
    this.messages = this.messagesSubject.asObservable();
  }

  public get messagesValue() {
    return this.messagesSubject.value;
  }

  getMessagesByConversationId(conversationId: number) {
    return this.http
      .get<Message[]>(
        `${environment.apiUrl}/v1/conversations/${conversationId}/messages`
      )
      .pipe(
        map((msgs) => {
          const currMsgs = { ...this.messagesValue };
          currMsgs[conversationId] = msgs;
          this.messagesSubject.next(currMsgs);
          return msgs;
        })
      );
  }

  getMessagesById(id: number) {
    return this.http
      .get<Message>(`${environment.apiUrl}/v1/messages/${id}`)
      .pipe(
        map((msg) => {
          this.addNextMessage(msg);
        })
      );
  }

  pushToLocalHistory(msg: Message) {
    const convExists =
      this.conversationService.conversationsValue[msg.conversationId];
    if (!convExists) {
      this.conversationService.getConversations().subscribe({
        error: (e) => console.log(e),
        next: (cs) => {
          this.conversationService.setUnreadMessage(msg);
        },
      });
    } else {
      this.conversationService.setUnreadMessage(msg);
    }
    const msgs = this.messagesValue[msg.conversationId];
    if (!msgs) {
      this.getMessagesByConversationId(msg.conversationId).subscribe({
        error: (e) => console.log(e),
        next: () => {
          this.addNextMessage(msg);
        },
      });
    } else {
      this.addNextMessage(msg);
    }
  }

  private addNextMessage(msg: Message) {
    const msgs = { ...this.messagesValue };
    const index = msgs[msg.conversationId].findIndex((m) => m.id === msg.id);
    if (index >= 0) {
      msgs[msg.conversationId][index] = msg;
    } else {
      msgs[msg.conversationId].push(msg);
    }
    this.messagesSubject.next(msgs);
  }
}
