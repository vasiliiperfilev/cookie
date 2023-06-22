import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormControl } from '@angular/forms';
import { User } from '@app/_models';
import { Conversation, ConversationDto } from '@app/_models/conversation';
import { ConversationsService, UserService } from '@app/_services';
import {
  debounceTime,
  distinctUntilChanged,
  filter,
  finalize,
  switchMap,
  tap,
} from 'rxjs';

@Component({
  selector: 'app-conversations',
  templateUrl: './conversations.component.html',
  styleUrls: ['./conversations.component.scss'],
})
export class ConversationsComponent implements OnInit {
  @Output() selectConversationEvent = new EventEmitter<Conversation>();
  loading = false;
  conversations: Conversation[] = [];
  user: User;
  userSearchControl = new FormControl('');
  searchedUsers: User[] = [];
  constructor(
    private conversationService: ConversationsService,
    private userService: UserService
  ) {
    this.user = userService.userValue!;
  }

  ngOnInit() {
    this.loading = true;
    this.conversationService
      .getConversations()
      .subscribe({ error: (err) => console.log(err) });
    this.conversationService.conversations.subscribe(
      (conversations) => (this.conversations = conversations)
    );
    this.autocomplete();
  }

  getInitials(name: string) {
    const initials = name.charAt(0) + name.charAt(1);
    return initials.toUpperCase();
  }

  addConversation() {
    const c: ConversationDto = {
      userIds: [3, 4],
    };
    this.conversationService
      .postConversation(c)
      .subscribe({ error: (err) => console.log(err) });
  }

  selectConversation(c: Conversation) {
    this.selectConversationEvent.emit(c);
  }

  autocomplete() {
    this.userSearchControl.valueChanges
      .pipe(
        filter((res) => {
          this.loading = true;
          this.searchedUsers = [];
          if (!res || res.length < 3) {
            this.loading = false;
            return false;
          }
          return true;
        }),
        distinctUntilChanged(),
        debounceTime(500),
        tap(() => {
          this.searchedUsers = [];
        }),
        switchMap((value) =>
          this.userService.getAllBySearch(value!).pipe(
            finalize(() => {
              this.loading = false;
            })
          )
        )
      )
      .subscribe((res: User[]) => {
        if (res) {
          this.searchedUsers = res.filter((u) => u.id !== this.user.id);
        } else {
          this.searchedUsers = [];
        }
      });
  }

  onUserSearchClick(id: number) {
    const existingConv = this.conversations.find((c) =>
      c.users.some((u) => u.id === id)
    );
    if (existingConv) {
      this.selectConversation(existingConv);
      this.userSearchControl.setValue('');
      return;
    }
    const c: ConversationDto = {
      userIds: [this.user.id, id],
    };
    this.conversationService.postConversation(c).subscribe({
      error: (err) => console.log(err),
      next: (c) => this.selectConversation(c),
    });
    this.userSearchControl.setValue('');
  }
}
