import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';

import {ApiService} from './api.service';

@Injectable()
export class ClipService {
  constructor (
    private apiService: ApiService,
  ) {}

  getClip(id: string) {
    return this.apiService.get(`/clips/${id}`);
  }

  addClip(id: string) {
    return this.apiService.post(`/clips/${id}`);
  }

  getTotalViews() {
    return this.apiService.get(`/clips/views`)
  }

  getTotalViewsPerStreamer() {
    return this.apiService.get(`/clips/views/perStreamer`)
  }
}
