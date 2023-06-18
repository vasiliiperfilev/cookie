import { Message, MessageDto } from './message';

export enum WsEventType {
  MESSAGE = 'message',
}

export interface WsMessageEvent {
  type: string;
  payload: Message | MessageDto;
}
