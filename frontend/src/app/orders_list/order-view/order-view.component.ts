import { Component, Inject, Optional } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Item, Order, OrderState } from '@app/_models';
import { ItemsService, UserService } from '@app/_services';
import { environment } from '@environments/environment';

interface OrderItem {
  item: Item;
  quantity: number;
}

@Component({
  selector: 'app-order-view',
  templateUrl: './order-view.component.html',
  styleUrls: ['./order-view.component.scss'],
})
export class OrderViewComponent {
  items: OrderItem[] = [];
  displayedColumns = ['image', 'name', 'size', 'unit', 'quantity'];
  public get OrderState() {
    return OrderState;
  }

  constructor(
    private itemsService: ItemsService,
    userService: UserService,
    public dialogRef: MatDialogRef<OrderViewComponent>,
    @Optional() @Inject(MAT_DIALOG_DATA) public data: Order
  ) {
    this.itemsService
      .getAllBySupplierId(userService.userValue!.id)
      .subscribe((items) => {
        this.items = items.map((item) => ({
          item,
          quantity: data.items.find((v) => v.itemId === item.id)!.quantity,
        }));
      });
  }

  onCancel(): void {
    this.dialogRef.close({});
  }

  getImageUrl(item: Item) {
    return `${environment.apiUrl}/v1/images/${item.imageId}`;
  }
}
