import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { Order, OrderState } from '@app/_models';
import { AlertService, ChatService, OrdersService } from '@app/_services';
import { CrudDialogAction } from '@app/catalog/catalog.component';
import { first } from 'rxjs';
import { OrderViewComponent } from './order-view/order-view.component';

@Component({
  selector: 'app-catalog',
  templateUrl: './order_list.component.html',
  styleUrls: ['./order_list.component.scss'],
})
export class OrderListComponent implements OnInit {
  @ViewChild('table') table: MatTable<any> | undefined;
  orders: Record<number, Order> = {};
  displayedColumns = ['id', 'client', 'createdAt', 'action'];
  public get CrudDialogAction() {
    return CrudDialogAction;
  }
  public get OrderState() {
    return OrderState;
  }
  public getOrders() {
    return Object.values(this.orders);
  }

  constructor(
    private orderService: OrdersService,
    private chatService: ChatService,
    private alertService: AlertService,
    public dialog: MatDialog
  ) {
    orderService.orders.subscribe((orders) => (this.orders = orders));
  }

  ngOnInit() {
    this.orderService.getAll().subscribe({
      error: (err) => console.log(err),
      next: (orders) => (this.orders = orders),
    });
  }

  canAcceptOrder(order: Order) {
    return (
      order.stateId !== OrderState.OrderStateAccepted &&
      order.stateId !== OrderState.OrderStateDeclined &&
      order.stateId !== OrderState.OrderStateFulfilled
    );
  }

  openOrderView(order: Order) {
    const dialogRef = this.dialog.open(OrderViewComponent, {
      width: '500px',
      data: order,
    });
  }

  updateOrder(orderId: number, stateId: OrderState) {
    this.orderService
      .update(orderId, {
        stateId,
      })
      .pipe(first())
      .subscribe({
        next: (order) => {
          this.alertService.success('Order state updated!');
          this.chatService.sendUpdatedOrder(order);
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
        },
      });
  }
}
