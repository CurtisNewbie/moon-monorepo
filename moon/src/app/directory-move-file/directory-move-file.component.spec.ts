import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DirectoryMoveFileComponent } from './directory-move-file.component';

describe('DirectoryMoveFileComponent', () => {
  let component: DirectoryMoveFileComponent;
  let fixture: ComponentFixture<DirectoryMoveFileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DirectoryMoveFileComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DirectoryMoveFileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
