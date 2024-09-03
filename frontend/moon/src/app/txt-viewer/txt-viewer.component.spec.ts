import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TxtViewerComponent } from './txt-viewer.component';

describe('TxtViewerComponent', () => {
  let component: TxtViewerComponent;
  let fixture: ComponentFixture<TxtViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TxtViewerComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TxtViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
