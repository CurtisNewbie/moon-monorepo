import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MngFilesComponent } from './mng-files.component';

describe('HomePageComponent', () => {
  let component: MngFilesComponent;
  let fixture: ComponentFixture<MngFilesComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ MngFilesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MngFilesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
