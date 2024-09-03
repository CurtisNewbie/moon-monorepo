import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MngResDialogComponent } from './mng-res-dialog.component';

describe('MngResDialogComponent', () => {
  let component: MngResDialogComponent;
  let fixture: ComponentFixture<MngResDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MngResDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MngResDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
