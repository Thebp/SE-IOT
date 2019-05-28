import { TestBed } from '@angular/core/testing';

import { WsadapterService } from './wsadapter.service';

describe('WsadapterService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: WsadapterService = TestBed.get(WsadapterService);
    expect(service).toBeTruthy();
  });
});
