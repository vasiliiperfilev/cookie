import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConversationsComponent } from './conversations.component';
import { MatListModule } from '@angular/material/list';
import { ConversationsService } from '@app/_services/conversations.service';

@NgModule({
  declarations: [ConversationsComponent],
  imports: [CommonModule, MatListModule],
  exports: [ConversationsComponent],
})
export class ConversationsModule {}
