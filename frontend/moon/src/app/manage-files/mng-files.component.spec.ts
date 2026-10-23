import { EventEmitter } from '@angular/core';
import { of, Subject, throwError } from 'rxjs';

import { FileInfo, FileType } from 'src/common/file-info';
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

/**
 * Stand-in for HttpClient that records what the component asked the server to
 * do. Each call answers with an empty, successful listing so the component can
 * run its normal fetch flow.
 */
class FakeHttp {
  posts: { url: string; body: any }[] = [];
  gets: string[] = [];
  /** value the order-by-preference endpoint answers with */
  orderByPreference: any = null;
  /** when set, the order-by-preference endpoint answers with an error */
  orderByPreferenceError = false;
  /** when set, the order-by-preference request fails outright */
  orderByPreferenceFails = false;

  post(url: string, body?: any): any {
    this.posts.push({ url, body });
    if (url.indexOf('dir/bottom-up-tree') >= 0) {
      return of({ data: null });
    }
    return of({ data: { payload: [], paging: { page: 1, limit: 30, total: 0 } } });
  }

  get(url: string): any {
    this.gets.push(url);
    if (url.indexOf('order-by-preference') >= 0) {
      if (this.orderByPreferenceFails) {
        return throwError(new Error('preference fetch failed'));
      }
      if (this.orderByPreferenceError) {
        return of({ error: true, msg: 'preference fetch failed' });
      }
      return of({ data: this.orderByPreference });
    }
    return of({ data: null });
  }

  postsTo(fragment: string): { url: string; body: any }[] {
    return this.posts.filter((p) => p.url.indexOf(fragment) >= 0);
  }

  lastPostTo(fragment: string): { url: string; body: any } {
    const matching = this.postsTo(fragment);
    return matching[matching.length - 1];
  }
}

describe('MngFilesComponent size sort', () => {
  let component: MngFilesComponent;
  let http: FakeHttp;
  let browseHistoryRecorder: {
    record: jasmine.Spy;
    recordDirPage: jasmine.Spy;
    getDirPage: jasmine.Spy;
  };
  let fileService: { generateFileTempToken: jasmine.Spy };
  /** route params the stubbed ActivatedRoute answers with */
  let routeParams: { [key: string]: string };

  beforeEach(() => {
    http = new FakeHttp();
    routeParams = {};
    browseHistoryRecorder = {
      record: jasmine.createSpy('record'),
      recordDirPage: jasmine.createSpy('recordDirPage'),
      getDirPage: jasmine.createSpy('getDirPage').and.returnValue(of('')),
    };
    fileService = {
      generateFileTempToken: jasmine
        .createSpy('generateFileTempToken')
        .and.returnValue(of({ error: true, msg: 'stubbed' })),
    };

    component = new MngFilesComponent(
      {} as any,
      { open: () => ({ afterClosed: () => of(false) }) } as any,
      fileService as any,
      { navigateTo: jasmine.createSpy('navigateTo') } as any,
      http as any,
      {
        paramMap: of({ get: (key: string) => routeParams[key] || null }),
      } as any,
      {} as any,
      { open: jasmine.createSpy('open') } as any,
      browseHistoryRecorder as any,
      {} as any,
      { isMobile: () => false } as any
    );

    component.pagingController = {
      paging: { page: 1, limit: 30, total: 0 },
      onTotalChanged: jasmine.createSpy('onTotalChanged'),
      atFirstPage: () => true,
      firstPage: jasmine.createSpy('firstPage'),
      syncPage: jasmine.createSpy('syncPage'),
    } as any;
  });

  function sortOptionValues(): string[] {
    return component.sortOptions.map((o) => o.value);
  }

  it('always offers size in the sort options', () => {
    component.searchParam = {};
    expect(sortOptionValues()).toEqual(['time', 'name', 'custom', 'size']);

    component.searchParam.fileType = FileType.DIR;
    expect(sortOptionValues()).toEqual(['time', 'name', 'custom', 'size']);

    component.searchParam.fileType = FileType.FILE;
    expect(sortOptionValues()).toEqual(['time', 'name', 'custom', 'size']);
  });

  it('keeps the sort options array stable across ticks and filter changes', () => {
    component.searchParam = { fileType: FileType.FILE };
    const atFile = component.sortOptions;

    expect(component.sortOptions).toBe(atFile);

    component.searchParam.fileType = FileType.DIR;
    expect(component.sortOptions).toBe(atFile);

    component.searchParam.fileType = FileType.FILE;
    expect(component.sortOptions).toBe(atFile);
    expect(sortOptionValues()).toEqual(['time', 'name', 'custom', 'size']);
  });

  it('choosing size under All narrows the listing to files in one fetch', () => {
    component.searchParam = {};
    component.orderBy = 'size';

    component.onOrderByChange();

    expect(component.searchParam.fileType).toBe(FileType.FILE);
    expect(http.postsTo('/file/list').length).toBe(1);
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.fileType).toBe(FileType.FILE);
    expect(listReq.orderBy).toBe('size');
    expect(listReq.orderDir).toBe('desc');
  });

  it('choosing size under Directory narrows the listing to files in one fetch', () => {
    component.searchParam.fileType = FileType.DIR;
    component.orderBy = 'size';

    component.onOrderByChange();

    expect(component.searchParam.fileType).toBe(FileType.FILE);
    expect(http.postsTo('/file/list').length).toBe(1);
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.fileType).toBe(FileType.FILE);
    expect(listReq.orderBy).toBe('size');
    expect(listReq.orderDir).toBe('desc');
  });

  it('choosing size while already on File keeps the filter', () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'size';

    component.onOrderByChange();

    expect(component.searchParam.fileType).toBe(FileType.FILE);
    expect(http.lastPostTo('/file/list').body.fileType).toBe(FileType.FILE);
  });

  it('leaving File while sorted by size restores the directory order', async () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'size';
    component.orderByDir = 'asc';
    http.orderByPreference = { orderBy: 'name', orderDir: 'asc' };

    component.searchParam.fileType = FileType.DIR;
    await component.onFileTypeChange();

    expect(component.orderBy).toBe('name');
    expect(component.orderByDir).toBe('asc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.fileType).toBe(FileType.DIR);
    expect(listReq.orderBy).toBe('name');
    expect(listReq.orderDir).toBe('asc');
  });

  it('leaving File while sorted by size falls back to the default order', async () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'size';

    component.searchParam.fileType = FileType.DIR;
    await component.onFileTypeChange();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('desc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.fileType).toBe(FileType.DIR);
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
  });

  it('leaving File while sorted by time keeps the chosen order', async () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'time';
    component.orderByDir = 'asc';
    http.orderByPreference = { orderBy: 'name', orderDir: 'asc' };

    component.searchParam.fileType = FileType.DIR;
    await component.onFileTypeChange();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('asc');
    expect(http.gets.some((u) => u.indexOf('order-by-preference') >= 0)).toBe(false);
  });

  it('does not upload a size sort as the order-by preference', () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'size';

    component.onOrderByChange();
    component.orderByDir = 'asc';
    component.onOrderByDirChange();

    expect(http.postsTo('order-by-preference')).toEqual([]);
    expect(http.lastPostTo('/file/list').body.orderDir).toBe('asc');
  });

  it('still uploads a durable sort as the order-by preference', () => {
    component.searchParam.fileType = FileType.FILE;
    component.orderBy = 'name';

    component.onOrderByChange();

    const pref = http.lastPostTo('order-by-preference').body;
    expect(pref.orderBy).toBe('name');
    expect(pref.orderDir).toBe('asc');
  });

  function flushAsync(): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve));
  }

  function previewFour(files: FileInfo[]): void {
    component.fileInfoList = files;
    for (let i = 0; i < files.length; i++) {
      component.preview(files[i], i);
    }
  }

  it('records no resume cursor while browsing a size-ordered listing', () => {
    component.inDirFileKey = 'dir-1';
    component.orderBy = 'size';

    previewFour([
      imgFile('u0', 'a.bin'),
      imgFile('u1', 'b.bin'),
      imgFile('u2', 'c.bin'),
      imgFile('u3', 'd.bin'),
    ]);

    expect(browseHistoryRecorder.recordDirPage).not.toHaveBeenCalled();
    // the plain browse history is still recorded, only the cursor is skipped
    expect(browseHistoryRecorder.record).toHaveBeenCalledTimes(4);
  });

  it('still records a resume cursor under a durable sort', () => {
    component.inDirFileKey = 'dir-1';
    component.orderBy = 'time';

    previewFour([
      imgFile('u0', 'a.bin'),
      imgFile('u1', 'b.bin'),
      imgFile('u2', 'c.bin'),
    ]);

    expect(browseHistoryRecorder.recordDirPage).toHaveBeenCalled();
  });

  it('does not offer to resume a directory under a size sort', async () => {
    routeParams = { parentDirKey: 'dir-1', orderBy: 'size' };

    component.parseRouteParam();
    await flushAsync();

    expect(browseHistoryRecorder.getDirPage).not.toHaveBeenCalled();
    expect(component.searchParam.fileType).toBe(FileType.FILE);
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('size');
    expect(listReq.fileType).toBe(FileType.FILE);
    expect(sortOptionValues()).toContain('size');
  });

  it('still offers to resume a directory under a durable sort', async () => {
    routeParams = { parentDirKey: 'dir-1', orderBy: 'time' };

    component.parseRouteParam();
    await flushAsync();

    expect(browseHistoryRecorder.getDirPage).toHaveBeenCalled();
  });

  it('restores the stored order when entering a directory that was sorted by size', async () => {
    component.orderBy = 'size';
    http.orderByPreference = { orderBy: 'name', orderDir: 'asc' };
    routeParams = { parentDirKey: 'dir-1' };

    component.parseRouteParam();
    await flushAsync();

    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('name');
    expect(listReq.orderDir).toBe('asc');
  });

  it('opens a directory in the default order when size was active and nothing is stored', async () => {
    component.orderBy = 'size';
    routeParams = { parentDirKey: 'dir-1' };

    component.parseRouteParam();
    await flushAsync();

    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
  });

  it('falls back to the default order when the preference fetch fails', async () => {
    component.orderBy = 'size';
    http.orderByPreferenceFails = true;
    routeParams = { parentDirKey: 'dir-1' };

    component.parseRouteParam();
    await flushAsync();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('desc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
  });

  it('falls back to the default order when the preference fetch returns an error', async () => {
    component.orderBy = 'size';
    http.orderByPreferenceError = true;
    routeParams = { parentDirKey: 'dir-1' };

    component.parseRouteParam();
    await flushAsync();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('desc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
  });

  it('never applies a stored size preference', async () => {
    component.orderBy = 'size';
    http.orderByPreference = { orderBy: 'size', orderDir: 'asc' };
    routeParams = { parentDirKey: 'dir-1' };

    component.parseRouteParam();
    await flushAsync();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('desc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
  });
});

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
