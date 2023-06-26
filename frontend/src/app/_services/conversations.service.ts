import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Conversation, ConversationDto } from '@app/_models/conversation';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { UserService } from './user.service';
import { Message } from '@app/_models';

@Injectable({
  providedIn: 'root',
})
export class ConversationsService {
  private conversationsSubject: BehaviorSubject<Record<number, Conversation>>;
  public conversations: Observable<Record<number, Conversation>>;

  constructor(private http: HttpClient, private userService: UserService) {
    this.conversationsSubject = new BehaviorSubject<
      Record<number, Conversation>
    >({});
    this.conversations = this.conversationsSubject.asObservable();
    this.getConversations().subscribe({
      error: (err) => console.log(err),
    });
  }

  public get conversationsValue() {
    return this.conversationsSubject.value;
  }

  getConversations() {
    return this.http
      .get<Conversation[]>(
        `${environment.apiUrl}/v1/conversations?userId=${this.userService.userValue?.id}&expanded=true`
      )
      .pipe(
        map((cs) => {
          const conversations: Record<number, Conversation> = cs.reduce(
            (acc, c) => {
              const conv = new Conversation(c.id, c.users, c.lastMessage);
              acc[c.id] = conv;
              return acc;
            },
            {} as Record<number, Conversation>
          );
          this.conversationsSubject.next(conversations);
          return conversations;
        })
      );
  }

  postConversation(cvs: ConversationDto) {
    return this.http
      .post<Conversation>(
        `${environment.apiUrl}/v1/conversations?userId=${this.userService.userValue?.id}`,
        cvs
      )
      .pipe(
        map((c) => {
          const newC = new Conversation(c.id, c.users, c.lastMessage);
          this.conversationsValue[newC.id] = newC;
          this.conversationsSubject.next(this.conversationsValue);
          return newC;
        })
      );
  }

  setUnreadMessage(msg: Message) {
    if (this.userService.userValue?.id !== msg.senderId) {
      this.conversationsValue[msg.conversationId].hasUnreadMsg = true;
      this.conversationsSubject.next(this.conversationsValue);
    }
  }

  setReadMessage(id: number) {
    this.conversationsValue[id].hasUnreadMsg = false;
    this.conversationsSubject.next(this.conversationsValue);
  }
}
