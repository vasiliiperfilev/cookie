export interface Order {
  id: number;
  messageId: number;
  createdAt: Date;
  updatedAt: Date;
  items: {
    itemId: number;
    quantity: number;
  }[];
  stateId: OrderState;
}

export enum OrderState {
  OrderStateCreated = 1,
  OrderStateAccepted = 2,
  OrderStateDeclined = 3,
  OrderStateFulfilled = 4,
  OrderStateConfirmedFulfillment = 5,
  OrderStateSupplierChanges = 6,
  OrderStateClientChanges = 7,
}

export interface PostOrderDto {
  items: {
    itemId: number;
    quantity: number;
  }[];
  conversationId: number;
}

export interface PatchOrderDto {
  items?: {
    itemId: number;
    quantity: number;
  }[];
  stateId?: OrderState;
}
