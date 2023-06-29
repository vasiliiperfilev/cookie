export interface Item {
  id: number;
  supplierId: number;
  unit: string;
  size: number;
  name: string;
  imageUrl: string;
}

export interface PostItemDto {
  unit: string;
  size: number;
  name: string;
  imageUrl: string;
}
