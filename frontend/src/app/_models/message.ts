import { User } from './user';

export interface Message {
  id: number;
  senderId: number;
  sender?: User;
  conversationId: number;
  prevMessageId: number;
  createdAt: Date;
  content: string;
}

export interface MessageDto {
  conversationId: number;
  prevMessageId: number;
  content: string;
}
