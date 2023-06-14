import { Message } from './message';

export enum WsEventType {
  MESSAGE = 'message',
}

export interface WsMessageEvent {
  type: string;
  payload: Message;
}
