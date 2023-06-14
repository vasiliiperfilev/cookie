import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { UserService } from './user.service';
import { Message } from '@app/_models/message';

@Injectable({
  providedIn: 'root',
})
export class HistoryService {
  private messagesSubject: BehaviorSubject<Message[]>;
  public messages: Observable<Message[]>;
  constructor(private http: HttpClient, private userService: UserService) {
    this.messagesSubject = new BehaviorSubject<Message[]>([]);
    this.messages = this.messagesSubject.asObservable();
  }

  public get messagesValue() {
    return this.messagesSubject.value;
  }

  getMessagesByConversationId(id: number) {
    return this.http
      .get<Message[]>(`${environment.apiUrl}/v1/conversations/${id}/messages`)
      .pipe(
        map((msgs) => {
          this.messagesSubject.next(msgs);
          return msgs;
        })
      );
  }

  pushToLocalHistory(msg: Message) {
    const msgs = [...this.messagesValue, msg];
    this.messagesSubject.next(msgs);
  }
}
