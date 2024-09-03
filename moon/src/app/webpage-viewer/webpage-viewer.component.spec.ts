import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WebpageViewerComponent } from './webpage-viewer.component';

describe('TxtViewerComponent', () => {
  let component: WebpageViewerComponent;
  let fixture: ComponentFixture<WebpageViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ WebpageViewerComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(WebpageViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
