import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConversationsComponent } from './conversations.component';
import { MatListModule } from '@angular/material/list';
import { MatToolbarModule } from '@angular/material/toolbar';

@NgModule({
  declarations: [ConversationsComponent],
  imports: [CommonModule, MatListModule, MatToolbarModule],
  exports: [ConversationsComponent],
})
export class ConversationsModule {}
