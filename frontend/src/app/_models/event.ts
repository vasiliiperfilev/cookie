import { Message, MessageDto } from './message';
import { Order } from './order';

export enum WsEventType {
  MESSAGE = 'message',
  ORDER = 'order',
}

export interface WsMessageEvent {
  type: string;
  payload: Message | MessageDto;
}

export interface WsOrderEvent {
  type: string;
  payload: Order;
}
