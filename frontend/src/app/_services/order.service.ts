import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Order, PatchOrderDto, PostOrderDto } from '@app/_models';
import { environment } from '@environments/environment';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class OrdersService {
  constructor(private http: HttpClient, private userService: UserService) {}

  getAll() {
    return this.http.get<Order[]>(
      `${environment.apiUrl}/v1/orders?userId=${this.userService.userValue?.id}`
    );
  }

  getById(id: number) {
    return this.http.get<Order>(`${environment.apiUrl}/v1/orders/${id}`);
  }

  create(dto: PostOrderDto) {
    return this.http.post<Order>(`${environment.apiUrl}/v1/orders`, dto);
  }

  update(id: number, dto: PatchOrderDto) {
    return this.http.patch<Order>(`${environment.apiUrl}/v1/orders/${id}`, dto);
  }
}
