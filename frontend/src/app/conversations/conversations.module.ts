import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConversationsComponent } from './conversations.component';
import { MatListModule } from '@angular/material/list';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AvatarComponent } from '@app/_components/avatar.component';
import { MatDividerModule } from '@angular/material/divider';

@NgModule({
  declarations: [ConversationsComponent],
  imports: [
    CommonModule,
    MatListModule,
    MatToolbarModule,
    AvatarComponent,
    MatDividerModule,
  ],
  exports: [ConversationsComponent],
})
export class ConversationsModule {}
