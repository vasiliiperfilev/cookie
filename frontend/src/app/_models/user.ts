export interface User {
  id: number;
  email: string;
  name: string;
  password?: string;
  type?: UserType;
  imageId: string;
}

export enum UserType {
  SUPPLIER = 1,
  BUSINESS = 2,
}

export interface PostUserDto {
  email: string;
  name: string;
  password: string;
  type: number;
  image?: File;
  imageId?: string;
}
