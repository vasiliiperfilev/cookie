export interface PostUserDto {
  email: string;
  name: string;
  password: string;
  type: number;
  imageId: string;
}

export interface PostTokenDto {
  email: string;
  password: string;
}
