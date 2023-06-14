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
