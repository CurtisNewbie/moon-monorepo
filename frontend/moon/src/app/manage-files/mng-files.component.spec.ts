import { EventEmitter } from '@angular/core';
import { of, Subject } from 'rxjs';

import { FileInfo } from 'src/common/file-info';
import { MngFilesComponent } from './mng-files.component';

/**
 * Controllable stand-in for the MatDialogRef returned by MatDialog.open().
 * afterOpened / afterClosed / keydownEvents are backed by Subjects the test
 * owns, so the test decides exactly when a transition completes.
 */
class FakeDialogRef {
  readonly componentInstance: any;
  readonly options: any;

  private afterOpened$ = new Subject<void>();
  private afterClosed$ = new Subject<void>();
  private keydown$ = new Subject<KeyboardEvent>();

  close = jasmine.createSpy('close').and.callFake(() => {
    this.afterClosed$.next();
    this.afterClosed$.complete();
  });

  constructor(componentInstance: any, options: any) {
    this.componentInstance = componentInstance;
    this.options = options;
  }

  afterOpened() {
    return this.afterOpened$.asObservable();
  }

  afterClosed() {
    return this.afterClosed$.asObservable();
  }

  keydownEvents() {
    return this.keydown$.asObservable();
  }

  fireAfterOpened() {
    this.afterOpened$.next();
  }
}

class FakeDialog {
  refs: FakeDialogRef[] = [];

  open(_component: any, options: any): FakeDialogRef {
    const componentInstance = {
      swipeLeft: new EventEmitter<boolean>(),
      swipeRight: new EventEmitter<boolean>(),
      prevEventEmitter: new EventEmitter<boolean>(),
      nextEventEmitter: new EventEmitter<boolean>(),
    };
    const ref = new FakeDialogRef(componentInstance, options);
    this.refs.push(ref);
    return ref;
  }
}

function imgFile(uuid: string, name: string): FileInfo {
  return { uuid, name, isFile: true } as FileInfo;
}

describe('MngFilesComponent preview nav coalescing', () => {
  let dialog: FakeDialog;
  let snackBar: { open: jasmine.Spy };
  let browseHistoryRecorder: { record: jasmine.Spy };
  let tokenResp: any;
  let fileService: { generateFileTempToken: jasmine.Spy };
  let component: MngFilesComponent;

  beforeEach(() => {
    dialog = new FakeDialog();
    snackBar = { open: jasmine.createSpy('open') };
    browseHistoryRecorder = { record: jasmine.createSpy('record') };
    tokenResp = { data: 'tok' };
    fileService = {
      generateFileTempToken: jasmine
        .createSpy('generateFileTempToken')
        .and.callFake(() => of(tokenResp)),
    };

    component = new MngFilesComponent(
      {} as any,
      dialog as any,
      fileService as any,
      { navigateTo: jasmine.createSpy('navigateTo') } as any,
      {} as any,
      {} as any,
      {} as any,
      snackBar as any,
      browseHistoryRecorder as any,
      {} as any,
      { isMobile: () => false } as any
    );

    component.fileInfoList = [
      imgFile('u0', 'img0.png'),
      imgFile('u1', 'img1.png'),
      imgFile('u2', 'img2.png'),
      imgFile('u3', 'img3.png'),
      imgFile('u4', 'img4.png'),
    ];
    (component as any).pagingController = {
      nextPage: () => false,
      prevPage: () => false,
    };
  });

  it('coalesces rapid swipes during a transition to the latest direction', () => {
    component.preview(component.fileInfoList[1], 1);
    expect(dialog.refs.length).toBe(1);
    expect(dialog.refs[0].options.data.name).toBe('img1.png');

    // Swipe to advance: starts a new (locked) transition to img2.
    dialog.refs[0].componentInstance.swipeRight.next(true);
    expect(dialog.refs.length).toBe(2);
    expect(dialog.refs[1].options.data.name).toBe('img2.png');

    // More swipes arrive while the new dialog is open but not yet shown.
    const pending = dialog.refs[1].componentInstance;
    pending.swipeLeft.next(true);
    pending.swipeRight.next(true);
    pending.swipeLeft.next(true);
    pending.swipeRight.next(true);
    pending.swipeLeft.next(true); // last one wins -> go left

    // The transition completes.
    dialog.refs[1].fireAfterOpened();

    // Exactly one coalesced navigation happened, in the last-swipe direction.
    expect(dialog.refs.length).toBe(3);
    expect(dialog.refs[2].options.data.name).toBe('img1.png');
  });

  it('does not drop a swipe that arrives while locked', () => {
    component.preview(component.fileInfoList[1], 1);

    dialog.refs[0].componentInstance.swipeRight.next(true); // -> img2, locked
    expect(dialog.refs.length).toBe(2);

    dialog.refs[1].componentInstance.swipeLeft.next(true); // single mid-transition swipe

    dialog.refs[1].fireAfterOpened();

    expect(dialog.refs.length).toBe(3);
    expect(dialog.refs[2].options.data.name).toBe('img1.png');
  });

  it('releases the nav lock when token generation fails', () => {
    component.preview(component.fileInfoList[1], 1);
    expect(dialog.refs.length).toBe(1);

    // The next transition's token generation fails.
    tokenResp = { error: true, msg: 'boom' };
    dialog.refs[0].componentInstance.swipeRight.next(true);

    expect(snackBar.open).toHaveBeenCalled();
    expect(dialog.refs.length).toBe(1); // no dialog opened for the failed target

    // The lock must be released: a later swipe still navigates.
    tokenResp = { data: 'tok' };
    component.preview(component.fileInfoList[2], 2);
    expect(dialog.refs.length).toBe(2);

    dialog.refs[1].componentInstance.swipeRight.next(true);

    expect(dialog.refs.length).toBe(3);
    expect(dialog.refs[2].options.data.name).toBe('img3.png');
  });
});
