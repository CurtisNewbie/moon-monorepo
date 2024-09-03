import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MngRoleDialogComponent } from './mng-role-dialog.component';

describe('MngRoleDialogComponent', () => {
  let component: MngRoleDialogComponent;
  let fixture: ComponentFixture<MngRoleDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MngRoleDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MngRoleDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
