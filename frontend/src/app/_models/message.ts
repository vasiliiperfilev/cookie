import { User } from './user';

export interface Message {
  id: number;
  sender: User;
  prevMessageId: number;
  createdAt: Date;
  content: string;
}
