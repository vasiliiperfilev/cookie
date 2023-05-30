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
