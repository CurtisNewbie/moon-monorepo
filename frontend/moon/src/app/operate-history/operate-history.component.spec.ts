import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OperateHistoryComponent } from './operate-history.component';

describe('OperateHistoryComponent', () => {
  let component: OperateHistoryComponent;
  let fixture: ComponentFixture<OperateHistoryComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ OperateHistoryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OperateHistoryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
