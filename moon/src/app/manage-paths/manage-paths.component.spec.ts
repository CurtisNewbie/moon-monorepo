import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ManagePathsComponent } from './manage-paths.component';

describe('ManagePathsComponent', () => {
  let component: ManagePathsComponent;
  let fixture: ComponentFixture<ManagePathsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ManagePathsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ManagePathsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
