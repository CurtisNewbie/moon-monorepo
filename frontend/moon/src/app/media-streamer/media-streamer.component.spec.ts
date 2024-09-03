import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MediaStreamerComponent } from './media-streamer.component';

describe('MediaStreamerComponent', () => {
  let component: MediaStreamerComponent;
  let fixture: ComponentFixture<MediaStreamerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MediaStreamerComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MediaStreamerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
