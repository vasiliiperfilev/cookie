import { HttpErrorResponse } from '@angular/common/http';
import { Component, Inject, Optional } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { FormErrors, PostItemDto } from '@app/_models';
import { AlertService } from '@app/_services';
import { ItemsService } from '@app/_services/items.service';
import { first } from 'rxjs';
import { CrudDialogAction, ItemDialogData } from '../catalog.component';

@Component({
  selector: 'app-create-item-dialog',
  templateUrl: './item-dialog.component.html',
  styleUrls: ['./item-dialog.component.scss'],
})
export class CreateItemDialogComponent {
  serverError: FormErrors<PostItemDto> | null = null;
  constructor(
    private alertService: AlertService,
    private itemService: ItemsService,
    public dialogRef: MatDialogRef<CreateItemDialogComponent>,
    @Optional() @Inject(MAT_DIALOG_DATA) public data: ItemDialogData
  ) {}

  form = new FormGroup({
    name: new FormControl(this.data.item?.name, {
      nonNullable: true,
      validators: [Validators.required],
    }),
    unit: new FormControl(this.data.item?.unit, {
      nonNullable: true,
      validators: [Validators.required],
    }),
    size: new FormControl(this.data.item?.size, {
      nonNullable: true,
      validators: [Validators.required],
    }),
    imageUrl: new FormControl('test', {
      nonNullable: true,
      validators: [Validators.required],
    }),
  });

  get f() {
    return this.form.controls;
  }

  getErrorMessage(fieldName: keyof PostItemDto) {
    if (this.form.controls[fieldName].hasError('required')) {
      return 'You must enter a value';
    } else if (this.serverError?.errors) {
      return this.serverError.errors[fieldName];
    }
    return '';
  }

  onSubmit() {
    this.alertService.clear();
    this.serverError = null;
    // stop here if form is invalid
    if (this.form.invalid) {
      return;
    }

    switch (this.data.action) {
      case CrudDialogAction.CREATE: {
        this.createItem();
        break;
      }
      case CrudDialogAction.UPDATE: {
        this.updateItem();
        break;
      }
      default: {
        this.deleteItem();
        break;
      }
    }
  }

  private deleteItem() {
    this.itemService
      .delete(this.data.item!.id)
      .pipe(first())
      .subscribe({
        next: () => {
          this.alertService.success('Item deleted!');
          this.dialogRef.close({
            action: CrudDialogAction.DELETE,
            item: this.data.item!,
          });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
          this.serverError = error.error as FormErrors<PostItemDto>;
        },
      });
  }

  private updateItem() {
    this.itemService
      .update(this.data.item!.id, {
        unit: this.form.value.unit!,
        size: this.form.value.size!,
        name: this.form.value.name!,
        imageUrl: this.form.value.imageUrl!,
      })
      .pipe(first())
      .subscribe({
        next: (item) => {
          this.alertService.success('Item updated!');
          this.dialogRef.close({ action: CrudDialogAction.UPDATE, item });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
          this.serverError = error.error as FormErrors<PostItemDto>;
        },
      });
  }

  private createItem() {
    this.itemService
      .create({
        unit: this.form.value.unit!,
        size: this.form.value.size!,
        name: this.form.value.name!,
        imageUrl: this.form.value.imageUrl!,
      })
      .pipe(first())
      .subscribe({
        next: (item) => {
          this.alertService.success('Added new item!');
          this.dialogRef.close({ action: CrudDialogAction.CREATE, item });
        },
        error: (error: HttpErrorResponse) => {
          this.alertService.error(error.error.message);
          console.log(error);
          this.serverError = error.error as FormErrors<PostItemDto>;
        },
      });
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
