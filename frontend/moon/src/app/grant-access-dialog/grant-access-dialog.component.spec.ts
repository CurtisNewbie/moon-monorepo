import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GrantAccessDialogComponent } from './grant-access-dialog.component';

describe('GrantAccessDialogComponent', () => {
  let component: GrantAccessDialogComponent;
  let fixture: ComponentFixture<GrantAccessDialogComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ GrantAccessDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GrantAccessDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
