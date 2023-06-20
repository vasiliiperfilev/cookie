import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConversationsComponent } from './conversations.component';
import { MatListModule } from '@angular/material/list';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AvatarComponent } from '@app/_components/avatar.component';

@NgModule({
  declarations: [ConversationsComponent],
  imports: [CommonModule, MatListModule, MatToolbarModule, AvatarComponent],
  exports: [ConversationsComponent],
})
export class ConversationsModule {}
