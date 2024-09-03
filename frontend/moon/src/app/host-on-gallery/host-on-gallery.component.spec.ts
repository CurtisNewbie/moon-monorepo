import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostOnGalleryComponent } from './host-on-gallery.component';

describe('HostOnGalleryComponent', () => {
  let component: HostOnGalleryComponent;
  let fixture: ComponentFixture<HostOnGalleryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HostOnGalleryComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HostOnGalleryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
