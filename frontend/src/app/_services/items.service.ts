import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Item, PostItemDto } from '@app/_models';
import { environment } from '@environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ItemsService {
  constructor(private http: HttpClient) {}

  getAll() {
    return this.http.get<Item[]>(`${environment.apiUrl}/v1/items`);
  }

  getById(id: number) {
    return this.http.get<Item>(`${environment.apiUrl}/v1/items/${id}`);
  }

  create(dto: PostItemDto) {
    return this.http.post(`${environment.apiUrl}/v1/items`, dto);
  }

  update(id: number, dto: PostItemDto) {
    return this.http.put(`${environment.apiUrl}/v1/items/${id}`, dto);
  }

  delete(id: number) {
    return this.http.delete(`${environment.apiUrl}/v1/items/${id}`);
  }
}
