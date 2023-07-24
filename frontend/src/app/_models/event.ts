import { Message, MessageDto } from './message';
import { Order } from './order';

export enum WsEventType {
  MESSAGE = 'message',
  NEW_ORDER = 'new_order',
  UPDATE_ORDER = 'update_order',
}

export interface WsMessageEvent {
  type: string;
  payload: Message | MessageDto;
}

export interface WsOrderEvent {
  type: string;
  payload: Order;
}
