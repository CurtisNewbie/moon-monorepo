import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GalleryAccessComponent } from './gallery-access.component';

describe('GalleryAccessComponent', () => {
  let component: GalleryAccessComponent;
  let fixture: ComponentFixture<GalleryAccessComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ GalleryAccessComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(GalleryAccessComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
