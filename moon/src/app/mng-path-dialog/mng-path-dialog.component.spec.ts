import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MngPathDialogComponent } from './mng-path-dialog.component';

describe('MngPathDialogComponent', () => {
  let component: MngPathDialogComponent;
  let fixture: ComponentFixture<MngPathDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MngPathDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MngPathDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
