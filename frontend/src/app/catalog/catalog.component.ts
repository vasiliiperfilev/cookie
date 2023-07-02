import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { Item } from '@app/_models';
import { ItemsService } from '@app/_services/items.service';
import { CreateItemDialogComponent } from './item-dialog/item-dialog.component';

export enum CrudDialogAction {
  CREATE = 'Create',
  UPDATE = 'Update',
  DELETE = 'Delete',
}
export interface ItemDialogData {
  action: CrudDialogAction;
  item?: Item;
}

@Component({
  selector: 'app-catalog',
  templateUrl: './catalog.component.html',
  styleUrls: ['./catalog.component.scss'],
})
export class CatalogComponent implements OnInit {
  @ViewChild('table') table: MatTable<any> | undefined;
  items: Item[] = [];
  selectedUnit: string | undefined;
  displayedColumns = ['name', 'size', 'unit', 'action'];
  public get CrudDialogAction() {
    return CrudDialogAction;
  }
  constructor(private itemService: ItemsService, public dialog: MatDialog) {}

  ngOnInit() {
    this.itemService.getAll().subscribe({
      error: (err) => console.log(err),
      next: (items) => (this.items = items),
    });
  }

  openDialog(itemDialogData: ItemDialogData): void {
    const dialogRef = this.dialog.open(CreateItemDialogComponent, {
      width: '250px',
      data: itemDialogData,
    });
    dialogRef.afterClosed().subscribe((result: ItemDialogData) => {
      if (result.item) {
        if (result.action == CrudDialogAction.CREATE) {
          this.addRowData(result.item);
        } else if (result.action == CrudDialogAction.UPDATE) {
          this.updateRowData(result.item);
        } else if (result.action == CrudDialogAction.DELETE) {
          this.deleteRowData(result.item);
        }
        this.table?.renderRows();
      }
    });
  }

  addRowData(item: Item) {
    this.items.push(item);
  }

  updateRowData(item: Item) {
    this.items = this.items.map((val) => {
      if (val.id == item.id) {
        return item;
      }
      return val;
    });
  }

  deleteRowData(item: Item) {
    this.items = this.items.filter((val) => {
      return val.id != item.id;
    });
  }
}
