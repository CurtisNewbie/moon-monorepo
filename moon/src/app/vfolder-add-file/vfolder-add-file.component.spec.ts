import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VfolderAddFileComponent } from './vfolder-add-file.component';

describe('VfolderAddFileComponent', () => {
  let component: VfolderAddFileComponent;
  let fixture: ComponentFixture<VfolderAddFileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VfolderAddFileComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VfolderAddFileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
