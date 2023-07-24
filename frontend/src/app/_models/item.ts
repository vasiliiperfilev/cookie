export interface Item {
  id: number;
  supplierId: number;
  unit: string;
  size: number;
  name: string;
  imageId: string;
}

export interface PostItemDto {
  unit: string;
  size: number;
  name: string;
  image?: File;
  imageId?: string;
}

export enum ItemUnit {
  l = 'liters',
  kg = 'kg',
}
