import { Message } from './message';

export interface Conversation {
  id: number;
  userIds: number[];
  lastMessage?: Message;
}

export interface ConversationDto {
  userIds: number[];
}
