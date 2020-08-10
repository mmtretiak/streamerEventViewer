import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';

import {ApiService} from './api.service';

@Injectable()
export class UserService {
  constructor (
    private apiService: ApiService,
  ) {}

  login() {
    return this.apiService.get('/users/login');
  }

  isLogged(): boolean {
    const token = localStorage.getItem("token");
    return token !== null;
  }

  getToken(code: string): Observable<any> {
    return this.apiService.get(`/users/login/redirect?code=${code}`);
  }
}
