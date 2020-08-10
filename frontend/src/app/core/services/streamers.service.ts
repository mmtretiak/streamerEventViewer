import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';

import {ApiService} from './api.service';

@Injectable()
export class StreamersService {
  constructor (
    private apiService: ApiService,
  ) {}

  getStreamers() {
    return this.apiService.get('/streamers');
  }

  addStreamer(name: string) {
    return this.apiService.post(`/streamers/${name}`);
  }
}
