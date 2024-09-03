import { TestBed } from '@angular/core/testing';

import { Toaster } from './notification.service';

describe('NotificationService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: Toaster = TestBed.get(Toaster);
    expect(service).toBeTruthy();
  });
});
