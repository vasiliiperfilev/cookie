import { HttpErrorResponse } from '@angular/common/http';
import { Component, Inject, Optional } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Item, UserType } from '@app/_models';
import { AlertService, ItemsService, OrdersService } from '@app/_services';
import { CrudDialogAction } from '@app/catalog/catalog.component';
import { first } from 'rxjs';
import { OrderDialogData } from '../chat.component';

@Component({
  selector: 'app-create-order-dialog',
  templateUrl: './order-dialog.component.html',
  styleUrls: ['./order-dialog.component.scss'],
})
export class OrderDialogComponent {
  items: Item[] = [];
  orderItems: { [x in number]: number } = {};
  displayedColumns = ['name', 'size', 'unit', 'quantity'];
  constructor(
    private alertService: AlertService,
    private orderService: OrdersService,
    private itemsService: ItemsService,
    public dialogRef: MatDialogRef<OrderDialogComponent>,
    @Optional() @Inject(MAT_DIALOG_DATA) public data: OrderDialogData
  ) {
    const supplier = data.conversation.users.find(
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
  }

  onSubmit() {
    this.alertService.clear();
    // stop here if form is invalid
    if (Object.values(this.orderItems).length === 0) {
      return;
    }

    switch (this.data.action) {
      case CrudDialogAction.CREATE: {
        this.createItem();
        break;
      }
      case CrudDialogAction.UPDATE: {
        this.updateOrder();
        break;
      }
    }
  }

  private updateOrder() {
    this.orderService
      .update(this.data.order!.id, {
        items: this.toItemsArray(),
      })
      .pipe(first())
      .subscribe({
        next: (order) => {
          this.alertService.success('Item updated!');
          this.dialogRef.close({ action: CrudDialogAction.UPDATE, order });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
        },
      });
  }

  private createItem() {
    this.orderService
      .create({
        conversationId: this.data.conversation.id,
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

  private toItemsArray() {
    const result: { itemId: number; quantity: number }[] = [];
    Object.keys(this.orderItems).map((k) => {
      result.push({
        itemId: Number.parseInt(k),
        quantity: this.orderItems[Number.parseInt(k)],
      });
    });
    return result;
  }
}
