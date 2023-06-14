import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Conversation } from '@app/_models/conversation';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class ConversationsService {
  private conversationsSubject: BehaviorSubject<Conversation[]>;
  public conversations: Observable<Conversation[]>;
  constructor(private http: HttpClient, private userService: UserService) {
    this.conversationsSubject = new BehaviorSubject<Conversation[]>([]);
    this.conversations = this.conversationsSubject.asObservable();
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
        map((conversations) => {
          this.conversationsSubject.next(conversations);
          return conversations;
        })
      );
  }

  postConversation(conversation: Conversation) {
    return this.http
      .post<Conversation>(
        `${environment.apiUrl}/v1/conversations?userId=${this.userService.userValue?.id}`,
        conversation
      )
      .pipe(
        map((conversation) => {
          const conversations = this.conversationsValue;
          this.conversationsSubject.next([...conversations, conversation]);
          return conversations;
        })
      );
  }

  /*
  Used to update conversation in a list after message was sent/received
  */
  updateConversation(conversation: Conversation) {
    const conversations = this.conversationsValue.map((curr) => {
      if (curr.id == conversation.id) {
        return conversation;
      }
      return curr;
    });
    this.conversationsSubject.next(conversations);
  }
}
