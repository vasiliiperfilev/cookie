import { HttpErrorResponse } from '@angular/common/http';
import { Component, Inject, Optional } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Item, OrderState, User, UserType } from '@app/_models';
import {
  AlertService,
  ItemsService,
  OrdersService,
  UserService,
} from '@app/_services';
import { CrudDialogAction } from '@app/catalog/catalog.component';
import { environment } from '@environments/environment';
import { first } from 'rxjs';
import { OrderDialogData } from '../chat.component';

@Component({
  selector: 'app-create-order-dialog',
  templateUrl: './order-dialog.component.html',
  styleUrls: ['./order-dialog.component.scss'],
})
export class OrderDialogComponent {
  items: Item[] = [];
  user: User;
  orderItems: { [x in number]: number } = {};
  displayedColumns = ['image', 'name', 'size', 'unit', 'quantity'];
  public get CrudDialogAction() {
    return CrudDialogAction;
  }
  public get OrderState() {
    return OrderState;
  }

  constructor(
    private alertService: AlertService,
    private orderService: OrdersService,
    private itemsService: ItemsService,
    private userService: UserService,
    public dialogRef: MatDialogRef<OrderDialogComponent>,
    @Optional() @Inject(MAT_DIALOG_DATA) public data: OrderDialogData
  ) {
    const supplier = data.conversation?.users.find(
      (u) => u.type === UserType.SUPPLIER
    );
    this.itemsService.getAllBySupplierId(supplier!.id).subscribe((items) => {
      this.items = items;
      items.map((v) => (this.orderItems[v.id] = 0));
    });
    if (data.order) {
      orderService.getById(data.order.id).subscribe((order) => {
        order.items.map((v) => (this.orderItems[v.itemId] = v.quantity));
      });
    }
    this.user = userService.userValue!;
    this.displayedColumns = [
      'image',
      'name',
      'size',
      'unit',
      data.action === CrudDialogAction.CREATE ? 'editQuantity' : 'quantity',
    ];
  }

  updateOrderState(state: OrderState) {
    this.orderService
      .update(this.data.order!.id, {
        stateId: state,
      })
      .pipe(first())
      .subscribe({
        next: (order) => {
          this.alertService.success('Order state updated!');
          this.dialogRef.close({ action: CrudDialogAction.UPDATE, order });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
        },
      });
  }

  createOrder() {
    this.orderService
      .create({
        conversationId: this.data.conversation!.id,
        items: this.toItemsArray(),
      })
      .pipe(first())
      .subscribe({
        next: (order) => {
          this.alertService.success('Added new order!');
          this.dialogRef.close({
            action: CrudDialogAction.CREATE,
            order,
          });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
        },
      });
  }

  onCancel(): void {
    this.dialogRef.close({});
  }

  public getOrderItems(items: Item[]) {
    return items.filter((i) => this.orderItems[i.id] > 0);
  }

  private toItemsArray() {
    const result: { itemId: number; quantity: number }[] = [];
    Object.keys(this.orderItems).map((k) => {
      if (this.orderItems[Number.parseInt(k)] > 0) {
        result.push({
          itemId: Number.parseInt(k),
          quantity: this.orderItems[Number.parseInt(k)],
        });
      }
    });
    return result;
  }

  getImageUrl(item: Item) {
    return `${environment.apiUrl}/v1/images/${item.imageId}`;
  }
}
