import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { UserService } from './user.service';
import { Message } from '@app/_models/message';
import { ConversationsService } from './conversations.service';
import { Conversation } from '@app/_models';

@Injectable({
  providedIn: 'root',
})
export class HistoryService {
  private messagesSubject: BehaviorSubject<Record<number, Message[]>>;
  public messages: Observable<Record<number, Message[]>>;
  constructor(
    private http: HttpClient,
    private conversationService: ConversationsService,
    private userService: UserService
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

  pushToLocalHistory(msg: Message) {
    const convExists = this.conversationService.conversationsValue.find(
      (c) => c.id === msg.conversationId
    );
    if (!convExists) {
      this.conversationService.getConversations().subscribe({
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
    msgs[msg.conversationId].push(msg);
    this.messagesSubject.next(msgs);
  }
}
