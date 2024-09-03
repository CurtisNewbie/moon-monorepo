import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BookmarkBlacklistComponent } from './bookmark-blacklist.component';

describe('BookmarkBlacklistComponent', () => {
  let component: BookmarkBlacklistComponent;
  let fixture: ComponentFixture<BookmarkBlacklistComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ BookmarkBlacklistComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(BookmarkBlacklistComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
