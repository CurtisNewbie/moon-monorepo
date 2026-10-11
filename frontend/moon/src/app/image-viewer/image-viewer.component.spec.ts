import { TestBed } from '@angular/core/testing';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import {
  Lightbox,
  LightboxConfig,
  LightboxEvent,
  LIGHTBOX_EVENT,
} from 'ngx-lightbox';
import { Subject } from 'rxjs';

import { Env } from 'src/common/env-util';
import {
  ImageViewerComponent,
  ImgViewerDialogData,
} from './image-viewer.component';

describe('ImageViewerComponent lightbox teardown', () => {
  let lightbox: { open: jasmine.Spy; close: jasmine.Spy };
  let config: any;
  let events: Subject<any>;
  let dialogRef: { close: jasmine.Spy };
  let data: ImgViewerDialogData;
  let component: ImageViewerComponent;

  beforeEach(() => {
    lightbox = {
      open: jasmine.createSpy('open'),
      close: jasmine.createSpy('close'),
    };
    config = {};
    events = new Subject<any>();
    dialogRef = { close: jasmine.createSpy('close') };
    data = { name: 'a.png', url: '/x', rotate: false };

    TestBed.configureTestingModule({
      providers: [
        { provide: Env, useValue: {} },
        { provide: Lightbox, useValue: lightbox },
        { provide: LightboxConfig, useValue: config },
        { provide: LightboxEvent, useValue: { lightboxEvent$: events } },
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    });

    component = new ImageViewerComponent(
      TestBed.inject(Env),
      TestBed.inject(Lightbox),
      TestBed.inject(LightboxConfig),
      TestBed.inject(LightboxEvent),
      TestBed.inject(MatDialogRef),
      TestBed.inject(MAT_DIALOG_DATA)
    );
  });

  it('closes the lightbox on destroy when no CLOSE event was emitted', () => {
    component.ngOnDestroy();

    expect(lightbox.close).toHaveBeenCalledTimes(1);
  });

  it('does not close the lightbox again when a CLOSE event was already emitted', () => {
    events.next({ id: LIGHTBOX_EVENT.CLOSE });
    expect(dialogRef.close).toHaveBeenCalledTimes(1);

    component.ngOnDestroy();

    expect(lightbox.close).not.toHaveBeenCalled();
  });
});
