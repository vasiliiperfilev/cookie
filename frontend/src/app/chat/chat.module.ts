import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatToolbarModule } from '@angular/material/toolbar';
import { AvatarComponent } from '@app/_components/avatar.component';
import { ChatComponent } from './chat.component';
import { ChatLayoutComponent } from './chat_layout.component';
import { ConversationsComponent } from './conversations.component';

@NgModule({
  declarations: [ChatComponent, ChatLayoutComponent, ConversationsComponent],
  imports: [
    CommonModule,
    MatInputModule,
    MatIconModule,
    MatCardModule,
    MatDividerModule,
    MatFormFieldModule,
    ReactiveFormsModule,
    MatGridListModule,
    MatListModule,
    MatButtonModule,
    MatToolbarModule,
    AvatarComponent,
    MatToolbarModule,
    MatAutocompleteModule,
  ],
  exports: [ChatComponent],
})
export class ChatModule {}
