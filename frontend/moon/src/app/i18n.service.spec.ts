import { TestBed } from '@angular/core/testing';

import { I18n } from './i18n.service';

describe('I18nService', () => {
  let service: I18n;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(I18n);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
