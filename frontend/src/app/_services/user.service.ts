import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { environment } from '@environments/environment';
import { User, UserRequest, UserResponse } from '@app/_models';
import { Token } from '@app/_models/token';

@Injectable({ providedIn: 'root' })
export class UserService {
  private userSubject: BehaviorSubject<User | null>;
  private tokenSubject: BehaviorSubject<Token | null>;
  public user: Observable<User | null>;

  constructor(private router: Router, private http: HttpClient) {
    this.userSubject = new BehaviorSubject(
      JSON.parse(localStorage.getItem('user')!)
    );
    this.tokenSubject = new BehaviorSubject(
      JSON.parse(localStorage.getItem('token')!)
    );
    this.user = this.userSubject.asObservable();
  }

  public get userValue() {
    return this.userSubject.value;
  }

  public get tokenValue() {
    return this.tokenSubject.value;
  }

  login(email: string, password: string) {
    return this.http
      .post<UserResponse>(`${environment.apiUrl}/token`, { email, password })
      .pipe(
        map((userResponse) => {
          // store jwt token in local storage to keep user logged in between page refreshes
          localStorage.setItem('token', JSON.stringify(userResponse.token));
          this.tokenSubject.next(userResponse.token);
          this.userSubject.next(userResponse.user);
          return userResponse;
        })
      );
  }

  logout() {
    // remove user from local storage and set current user to null
    localStorage.removeItem('token');
    this.tokenSubject.next(null);
    this.userSubject.next(null);
    this.router.navigate(['/account/login']);
  }

  register(user: UserRequest) {
    return this.http.post(`${environment.apiUrl}/users`, user);
  }

  getById(id: string) {
    return this.http.get<User>(`${environment.apiUrl}/user/${id}`);
  }

  update(user: User) {
    return this.http
      .put<User>(`${environment.apiUrl}/user/${user.id}`, user)
      .pipe(
        map((user) => {
          // publish updated user to subscribers
          this.userSubject.next(user);
          return user;
        })
      );
  }

  delete(id: string) {
    return this.http.delete(`${environment.apiUrl}/user/${id}`).pipe(
      map((x) => {
        // auto logout if the logged in user deleted their own record
        if (id == this.userValue?.id) {
          this.logout();
        }
        return x;
      })
    );
  }
}