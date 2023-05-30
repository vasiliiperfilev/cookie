export class User {
  id?: string;
  email?: string;
  password?: string;
  type?: number;
  imageId?: string;
}

export enum UserType {
  SUPPLIER = 1,
  BUSINESS = 2,
}
