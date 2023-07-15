import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Order, PatchOrderDto, PostOrderDto } from '@app/_models';
import { environment } from '@environments/environment';
import { BehaviorSubject, Observable, map } from 'rxjs';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class OrdersService {
  private ordersSubject: BehaviorSubject<Record<number, Order>>;
  public orders: Observable<Record<number, Order>>;

  constructor(private http: HttpClient, private userService: UserService) {
    this.ordersSubject = new BehaviorSubject<Record<number, Order>>({});
    this.orders = this.ordersSubject.asObservable();
  }

  getAll() {
    return this.http
      .get<Order[]>(
        `${environment.apiUrl}/v1/orders?userId=${this.userService.userValue?.id}`
      )
      .pipe(
        map((os) => {
          const conversations: Record<number, Order> = os.reduce((acc, o) => {
            acc[o.messageId] = o;
            return acc;
          }, {} as Record<number, Order>);
          this.ordersSubject.next(conversations);
          return conversations;
        })
      );
  }

  getById(id: number) {
    return this.http.get<Order>(`${environment.apiUrl}/v1/orders/${id}`).pipe(
      map((o) => {
        return this.addOrderToObservable(o);
      })
    );
  }

  create(dto: PostOrderDto) {
    return this.http.post<Order>(`${environment.apiUrl}/v1/orders`, dto).pipe(
      map((o) => {
        return this.addOrderToObservable(o);
      })
    );
  }

  update(id: number, dto: PatchOrderDto) {
    return this.http
      .patch<Order>(`${environment.apiUrl}/v1/orders/${id}`, dto)
      .pipe(
        map((o) => {
          return this.addOrderToObservable(o);
        })
      );
  }

  private addOrderToObservable(o: Order) {
    const orders = { ...this.ordersSubject.value };
    orders[o.messageId] = o;
    this.ordersSubject.next(orders);
    return o;
  }

  pushToLocalOrders(order: Order) {
    this.addOrderToObservable(order);
  }
}
