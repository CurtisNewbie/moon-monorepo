import { TestBed } from '@angular/core/testing';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import {
  MediaStreamerComponent,
  MediaStreamerDialogData,
} from './media-streamer.component';

describe('MediaStreamerComponent', () => {
  let data: MediaStreamerDialogData;
  let component: MediaStreamerComponent;

  beforeEach(() => {
    data = {
      name: 'v.mp4',
      url: 'fstore/file/stream?key=tok',
      token: 'tok',
    };

    TestBed.configureTestingModule({
      providers: [
        { provide: MatDialogRef, useValue: { close: () => undefined } },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    });

    component = new MediaStreamerComponent(
      TestBed.inject(MatDialogRef),
      TestBed.inject(MAT_DIALOG_DATA)
    );
    component.ngOnInit();
  });

  it('pauses and clears the video source on destroy', () => {
    const fakeVideo = {
      pause: jasmine.createSpy('pause'),
      removeAttribute: jasmine.createSpy('removeAttribute'),
      load: jasmine.createSpy('load'),
    };
    component.videoPlayer = { nativeElement: fakeVideo } as any;

    component.ngOnDestroy();

    expect(fakeVideo.pause).toHaveBeenCalled();
    expect(fakeVideo.removeAttribute).toHaveBeenCalledWith('src');
    expect(fakeVideo.load).toHaveBeenCalled();
  });

  it('builds the stream url from the current origin', () => {
    expect(component.srcUrl).toBe(
      location.origin + '/fstore/file/stream?key=tok'
    );
  });

  it('marks the preview as failed when the video errors', () => {
    expect(component.loadFailed).toBe(false);

    component.onVideoError();

    expect(component.loadFailed).toBe(true);
  });
});
