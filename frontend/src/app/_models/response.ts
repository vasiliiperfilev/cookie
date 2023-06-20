import { PostUserDto } from './request';
import { Token } from './token';
import { User } from './user';

export class UserResponse {
  user: User;
  token: Token;
  constructor(user: User, token: Token) {
    this.user = user;
    this.token = token;
  }
}

export class FormErrors<T> {
  message: string;
  errors: Record<keyof T, string>;

  constructor(message: string, error: Record<keyof T, string>) {
    this.errors = error;
    this.message = message;
  }
}
