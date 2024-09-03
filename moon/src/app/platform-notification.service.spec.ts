import { TestBed } from '@angular/core/testing';

import { PlatformNotificationService } from './platform-notification.service';

describe('PlatformNotificationService', () => {
  let service: PlatformNotificationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PlatformNotificationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
