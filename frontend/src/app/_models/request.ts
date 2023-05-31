export class UserRequest {
  email: string;
  password: string;
  type: number;
  imageId: string;
  constructor(email: string, password: string, type: number, imageId: string) {
    this.email = email;
    this.password = password;
    this.type = type;
    this.imageId = imageId;
  }
}

export class TokenRequest {
  email: string;
  password: string;
  constructor(email: string, password: string) {
    this.email = email;
    this.password = password;
  }
}
