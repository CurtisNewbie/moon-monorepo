import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ControlledPaginatorComponent } from './controlled-paginator.component';

describe('ControlledPaginatorComponent', () => {
  let component: ControlledPaginatorComponent;
  let fixture: ComponentFixture<ControlledPaginatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ControlledPaginatorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ControlledPaginatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
