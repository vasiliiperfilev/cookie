import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'avatar',
  templateUrl: 'avatar.component.html',
  styleUrls: ['./avatar.component.scss'],
  standalone: true,
  imports: [CommonModule],
})
export class AvatarComponent {
  @Input({ required: true }) fallback!: string;
  @Input() imageUrl: string | undefined;
  getInitials() {
    const initials = this.fallback.charAt(0) + this.fallback.charAt(1);
    return initials.toUpperCase();
  }
}
