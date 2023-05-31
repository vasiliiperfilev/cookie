import { UserRequest } from './request';
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

export class FormErrors {
  error: Record<keyof UserRequest, string>;

  constructor(error: Record<string, string>) {
    this.error = error;
  }
}
