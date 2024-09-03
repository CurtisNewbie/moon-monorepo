import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManageLogsComponent } from './manage-logs.component';

describe('ManageLogsComponent', () => {
  let component: ManageLogsComponent;
  let fixture: ComponentFixture<ManageLogsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ManageLogsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManageLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
