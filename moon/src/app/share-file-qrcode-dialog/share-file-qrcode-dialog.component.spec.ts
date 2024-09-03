import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShareFileQrcodeDialogComponent } from './share-file-qrcode-dialog.component';

describe('ShareFileQrcodeDialogComponent', () => {
  let component: ShareFileQrcodeDialogComponent;
  let fixture: ComponentFixture<ShareFileQrcodeDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ShareFileQrcodeDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ShareFileQrcodeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
