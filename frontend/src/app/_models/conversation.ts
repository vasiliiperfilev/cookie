import { environment } from '@environments/environment';
import { Message } from './message';
import { User } from './user';

export class Conversation {
  id: number;
  users: User[];
  lastMessage?: Message;
  hasUnreadMsg = false;
  constructor(id: number, users: User[], lastMessage?: Message) {
    this.id = id;
    this.users = users;
    this.lastMessage = lastMessage;
  }

  getName(userId: number) {
    if (this.users.length == 2) {
      const user = this.users.find((user) => user.id !== userId);
      return user?.name ?? null;
    }
    return null;
  }

  getImageUrl(userId: number) {
    if (this.users.length == 2) {
      const user = this.users.find((user) => user.id !== userId);
      if (user) {
        return `${environment.apiUrl}/v1/images/${user.imageId}`;
      }
    }
    return undefined;
  }
}

export interface ConversationDto {
  userIds: number[];
}
