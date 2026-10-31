import { CommonModule, Location } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { EventEmitter, Component, Input, Output } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute } from '@angular/router';
import { BrowseHistoryRecorder } from 'src/common/browse-history';
import { FileBookmark } from 'src/common/file-bookmark';
import { I18n } from '../i18n.service';
import { NavigationService } from '../navigation.service';
import { FormsModule } from '@angular/forms';
import { MatBadgeModule } from '@angular/material/badge';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { DomSanitizer } from '@angular/platform-browser';
import { BehaviorSubject, of, Subject, throwError } from 'rxjs';

import { Env } from 'src/common/env-util';
import { FileInfo, FileType } from 'src/common/file-info';
import {
  FakeMediaQueryList,
  installFakeMatchMedia,
} from 'src/testing/fake-media-query-list';
import { FileInfoService } from '../file-info.service';
import { ImageTooltipDirective } from '../shared/directives/image-tooltip.directive';
import { TrlPipe } from '../trl.pipe';
import { MngFilesComponent } from './mng-files.component';

/** Placeholder for the paginator child: the folder preview tests render no rows through it. */
@Component({ selector: 'app-controlled-paginator', template: '' })
class StubPaginatorComponent {
  @Input() showGotoPageOnMobile: any;
  @Output() pageChanged = new EventEmitter<any>();
  @Output() ready = new EventEmitter<void>();
}

/** Placeholder for the file tray child: it is only rendered when the tray opens. */
@Component({ selector: 'app-file-tray', template: '' })
class StubFileTrayComponent {
  @Output() closed = new EventEmitter<void>();
}

/** Layout signal for tests that do not care which mode is active. */
function staticEnv(isMobile: boolean): Env {
  return { isMobile: () => isMobile, mobile$: of(isMobile) } as any;
}

/** Lets the zone settle the async work a test started. */
function flushAsync(): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve));
}

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
  /** value the directory zip endpoint answers with */
  dirZipResp: any = { data: 'bundle-key' };
  /** when set, the directory zip request fails at the transport level */
  dirZipFails = false;
  /** when set, the directory zip request stays pending until the test answers it */
  dirZipPending: Subject<any> = null;

  post(url: string, body?: any): any {
    this.posts.push({ url, body });
    if (url.indexOf('dir/bottom-up-tree') >= 0) {
      return of({ data: null });
    }
    if (url.indexOf('dir/zip') >= 0) {
      if (this.dirZipPending) {
        return this.dirZipPending.asObservable();
      }
      if (this.dirZipFails) {
        return throwError(new Error('zip request failed'));
      }
      return of(this.dirZipResp);
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
      staticEnv(false)
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

  it('reset while sorted by size returns to time descending', () => {
    component.searchParam = { fileType: FileType.FILE, name: 'report' };
    component.orderBy = 'size';
    component.orderByDir = 'asc';

    component.resetSearchParam();

    expect(component.orderBy).toBe('time');
    expect(component.orderByDir).toBe('desc');
    expect(component.searchParam.fileType).toBeUndefined();
    expect(component.searchParam.name).toBeUndefined();
    expect(http.postsTo('order-by-preference')).toEqual([]);
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('time');
    expect(listReq.orderDir).toBe('desc');
    expect(listReq.fileType).toBeUndefined();
    expect(listReq.filename).toBeUndefined();
  });

  it('reset keeps a non-size sort as it is', () => {
    component.searchParam = { fileType: FileType.FILE, name: 'report' };
    component.orderBy = 'name';
    component.orderByDir = 'asc';

    component.resetSearchParam();

    expect(component.orderBy).toBe('name');
    expect(component.orderByDir).toBe('asc');
    const listReq = http.lastPostTo('/file/list').body;
    expect(listReq.orderBy).toBe('name');
    expect(listReq.orderDir).toBe('asc');
    expect(listReq.fileType).toBeUndefined();
  });

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
      imgFile('u3', 'd.bin'),
      imgFile('u4', 'e.bin'),
    ]);

    expect(browseHistoryRecorder.recordDirPage).toHaveBeenCalled();
  });

  it('waits for a 5-view streak before recording a resume cursor', () => {
    component.inDirFileKey = 'dir-1';
    component.orderBy = 'time';

    previewFour([
      imgFile('u0', 'a.bin'),
      imgFile('u1', 'b.bin'),
      imgFile('u2', 'c.bin'),
      imgFile('u3', 'd.bin'),
    ]);

    expect(browseHistoryRecorder.recordDirPage).not.toHaveBeenCalled();
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
      staticEnv(false)
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

describe('MngFilesComponent directory zip download', () => {
  let component: MngFilesComponent;
  let http: FakeHttp;
  let snackBar: { open: jasmine.Spy };
  let openSpy: jasmine.Spy;

  beforeEach(() => {
    http = new FakeHttp();
    snackBar = { open: jasmine.createSpy('open') };
    openSpy = spyOn(window, 'open').and.stub();

    const fileService = new FileInfoService(http as any, snackBar as any);
    component = new MngFilesComponent(
      {} as any,
      {} as any,
      fileService,
      {} as any,
      http as any,
      {} as any,
      {} as any,
      snackBar as any,
      {} as any,
      {} as any,
      staticEnv(false)
    );
  });

  function dir(name: string): FileInfo {
    return { uuid: 'dir-1', name, isDir: true } as FileInfo;
  }

  it('posts the directory file key and opens the archive built from the returned key', () => {
    http.dirZipResp = { data: 'bundle key/1' };

    component.downloadDirAsZip(dir('Photos'));

    expect(http.lastPostTo('/file/dir/zip').body).toEqual({ fileKey: 'dir-1' });
    expect(openSpy).toHaveBeenCalledWith(
      'fstore/file/bundle?key=bundle%20key%2F1',
      '_parent'
    );
  });

  it('shows the server message instead of opening an archive when the key request is refused', () => {
    http.dirZipResp = { error: true, msg: 'not permitted' };

    component.downloadDirAsZip(dir('Photos'));

    expect(openSpy).not.toHaveBeenCalled();
    expect(snackBar.open).toHaveBeenCalledWith('not permitted', 'ok', {
      duration: 6000,
    });
  });

  it('collapses the expanded panel once the archive key arrives', () => {
    component.curr = dir('Photos');

    component.downloadDirAsZip(dir('Photos'));

    expect(component.curr).toBeNull();
  });

  it('keeps the panel open when the key request is refused', () => {
    http.dirZipResp = { error: true, msg: 'not permitted' };
    const d = dir('Photos');
    component.curr = d;

    component.downloadDirAsZip(d);

    expect(component.curr).toBe(d);
  });

  it('does nothing for a file row', () => {
    component.downloadDirAsZip(imgFile('u0', 'a.png'));

    expect(http.postsTo('/file/dir/zip')).toEqual([]);
    expect(openSpy).not.toHaveBeenCalled();
  });

  it('reports a transport failure and opens no archive', () => {
    http.dirZipFails = true;

    component.downloadDirAsZip(dir('Photos'));

    expect(openSpy).not.toHaveBeenCalled();
    expect(snackBar.open).toHaveBeenCalledWith(
      'Request failed, unknown error',
      'ok',
      { duration: 3000 }
    );
  });

  it('lets the user retry after a transport failure', () => {
    http.dirZipFails = true;
    component.downloadDirAsZip(dir('Photos'));

    http.dirZipFails = false;
    component.downloadDirAsZip(dir('Photos'));

    expect(http.postsTo('/file/dir/zip').length).toBe(2);
    expect(openSpy).toHaveBeenCalledWith(
      'fstore/file/bundle?key=bundle-key',
      '_parent'
    );
  });

  it('asks for one key while a request is already in flight', () => {
    http.dirZipPending = new Subject<any>();

    component.downloadDirAsZip(dir('Photos'));
    component.downloadDirAsZip(dir('Photos'));

    expect(http.postsTo('/file/dir/zip').length).toBe(1);

    http.dirZipPending.next({ data: 'bundle-key' });
    http.dirZipPending.complete();

    expect(openSpy).toHaveBeenCalledTimes(1);
    expect(openSpy).toHaveBeenCalledWith(
      'fstore/file/bundle?key=bundle-key',
      '_parent'
    );
  });

  it('asks for a new key once the previous request finished', () => {
    component.downloadDirAsZip(dir('Photos'));
    component.downloadDirAsZip(dir('Photos'));

    expect(http.postsTo('/file/dir/zip').length).toBe(2);
    expect(openSpy).toHaveBeenCalledTimes(2);
  });
});

describe('MngFilesComponent columns', () => {
  let fake: FakeMediaQueryList;
  let component: MngFilesComponent;

  const WIDE = ['thumbnail', 'name', 'parentFileName', 'uploadTime', 'size', 'operation'];
  const COMPACT = ['thumbnail', 'name', 'operation'];

  beforeEach(() => {
    fake = installFakeMatchMedia(false);
    component = new MngFilesComponent(
      {} as any,
      { open: () => ({ afterClosed: () => of(false) }) } as any,
      {} as any,
      { navigateTo: jasmine.createSpy('navigateTo') } as any,
      {} as any,
      { paramMap: of({ get: () => null }) } as any,
      {} as any,
      { open: jasmine.createSpy('open') } as any,
      { record: jasmine.createSpy('record') } as any,
      {} as any,
      new Env()
    );
  });

  it('shows the wide columns at 1280px and the compact columns at 390px', () => {
    expect(component.displayedColumns).toEqual(WIDE);

    fake.fireChange(true);
    expect(component.displayedColumns).toEqual(COMPACT);
  });

  it('reassigns the columns when the layout flips mid-life', () => {
    fake.fireChange(true);
    expect(component.displayedColumns).toEqual(COMPACT);

    fake.fireChange(false);
    expect(component.displayedColumns).toEqual(WIDE);
  });

  it('keeps the drag handle out of the compact list', () => {
    component.orderBy = 'custom';
    expect(component.displayedColumns[0]).toBe('dragHandle');

    fake.fireChange(true);
    expect(component.displayedColumns).toEqual(COMPACT);
  });
});

/**
 * Folder preview, rendered through the real manage-files template with the
 * HTTP client stubbed. The assertions are what the user sees in the Preview
 * cell, plus what the browser remembers after the page is reopened.
 */
describe('MngFilesComponent folder preview', () => {
  let fixture: ComponentFixture<MngFilesComponent>;
  let component: MngFilesComponent;
  let http: FakeHttp;
  let mobile$: BehaviorSubject<boolean>;

  /** a real, decodable image so intrinsic thumbnail sizes are measurable */
  const IMAGE_DATA_URL =
    'data:image/svg+xml;charset=utf-8,' +
    encodeURIComponent('<svg xmlns="http://www.w3.org/2000/svg" width="160" height="100"></svg>');
  const COVER_URL = IMAGE_DATA_URL;
  const FOLDER_PREVIEW_KEY = 'manage-files.folderPreview';

  function dirRow(uuid: string, name: string): FileInfo {
    return { uuid, name, isDir: true, fileType: 'DIR' } as FileInfo;
  }

  /** A data URL the browser will actually decode; Angular would prefix it as unsafe otherwise. */
  function trustedImage(url: string): any {
    return TestBed.inject(DomSanitizer).bypassSecurityTrustUrl(url);
  }

  function fileRow(uuid: string, name: string, thumbnailUrl: string): FileInfo {
    return { uuid, name, isFile: true, fileType: 'FILE', thumbnailUrl } as FileInfo;
  }

  function open(rows: FileInfo[]): void {
    fixture = TestBed.createComponent(MngFilesComponent);
    component = fixture.componentInstance;
    component.fileInfoList = rows;
    fixture.detectChanges();
  }

  function checkbox(label: string, root: HTMLElement = fixture.nativeElement): HTMLElement {
    const boxes = Array.from<HTMLElement>(root.querySelectorAll('mat-checkbox'));
    return boxes.find((box) => (box.textContent || '').indexOf(label) >= 0);
  }

  function isChecked(label: string, root: HTMLElement = fixture.nativeElement): boolean {
    return (checkbox(label, root).querySelector('input') as HTMLInputElement).checked;
  }

  function tickFolderPreview(): void {
    (checkbox('folderPreview').querySelector('input') as HTMLInputElement).click();
    fixture.detectChanges();
  }

  function rowAt(index: number): HTMLElement {
    return fixture.nativeElement.querySelectorAll('tr.file-row')[index];
  }

  function cellAt(index: number): HTMLElement {
    return rowAt(index).querySelector('td.mat-column-thumbnail');
  }

  function thumbnailCell(): HTMLElement {
    return cellAt(0);
  }

  function glyphAt(index: number): HTMLElement {
    return cellAt(index).querySelector('.dir-glyph');
  }

  /** Lets real image decoding finish before the test measures a box. */
  async function settleImages(): Promise<void> {
    fixture.detectChanges();
    const imgs = Array.from<HTMLImageElement>(
      fixture.nativeElement.querySelectorAll('img')
    );
    await Promise.all(
      imgs.map((img) =>
        img.complete
          ? Promise.resolve()
          : new Promise<void>((resolve) => {
              img.addEventListener('load', () => resolve(), { once: true });
              img.addEventListener('error', () => resolve(), { once: true });
            })
      )
    );
    fixture.detectChanges();
  }

  /**
   * Renders at a given width for the duration of `run`. The layout tiers live
   * in media queries and Karma's own window is 756px wide, so the runner iframe
   * is resized — the assertion is still layout a real browser did.
   */
  function atViewportWidth(
    width: number,
    run: (viewportWidth: number) => void
  ): void {
    const frame = window.frameElement as HTMLElement;
    const originalWidth = frame ? frame.style.width : null;
    if (frame) {
      frame.style.width = `${width}px`;
    }
    try {
      const viewportWidth = frame
        ? frame.getBoundingClientRect().width
        : document.body.getBoundingClientRect().width;
      run(viewportWidth);
    } finally {
      if (frame) {
        frame.style.width = originalWidth;
      }
    }
  }

  function atPhoneWidth(run: (viewportWidth: number) => void): void {
    atViewportWidth(412, run);
  }

  beforeEach(() => {
    // a fresh browser for this page's preference only
    localStorage.removeItem(FOLDER_PREVIEW_KEY);
    http = new FakeHttp();
    mobile$ = new BehaviorSubject<boolean>(false);

    TestBed.configureTestingModule({
      declarations: [
        MngFilesComponent,
        TrlPipe,
        ImageTooltipDirective,
        StubPaginatorComponent,
        StubFileTrayComponent,
      ],
      imports: [
        CommonModule,
        FormsModule,
        NoopAnimationsModule,
        MatBadgeModule,
        MatButtonModule,
        MatCardModule,
        MatCheckboxModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatSelectModule,
        MatTableModule,
        MatTooltipModule,
        DragDropModule,
      ],
      providers: [
        { provide: I18n, useValue: { trl: (mod: string, key: string) => key } },
        { provide: MatDialog, useValue: new FakeDialog() },
        {
          provide: FileInfoService,
          useValue: {
            generateFileTempToken: () => of({ error: true, msg: 'stubbed' }),
          },
        },
        { provide: NavigationService, useValue: { navigateTo: () => undefined } },
        { provide: HttpClient, useValue: http },
        {
          provide: ActivatedRoute,
          useValue: { paramMap: of({ get: () => null }) },
        },
        { provide: Location, useValue: {} },
        { provide: MatSnackBar, useValue: { open: () => undefined } },
        { provide: BrowseHistoryRecorder, useValue: { record: () => undefined } },
        { provide: FileBookmark, useValue: { isEmpty: () => true, has: () => false } },
        {
          provide: Env,
          useValue: {
            isMobile: () => mobile$.value,
            mobile$: mobile$.asObservable(),
          },
        },
      ],
    });
  });

  it('starts off and draws the empty folder shape for a directory with a cached cover', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    fixture.detectChanges();

    expect(isChecked('folderPreview')).toBe(false);
    const cell = thumbnailCell();
    const glyph = glyphAt(0);
    expect(glyph).not.toBeNull();
    expect(glyph.querySelector('img')).toBeNull();
    // the shape itself is the surface: no inner rectangle and no bootstrap glyph
    expect(glyph.children.length).toBe(0);
    expect(cell.querySelector('i.bi')).toBeNull();
  });

  it('remembers a ticked box when the page is opened again', async () => {
    open([dirRow('dir-1', 'Album')]);
    expect(isChecked('folderPreview')).toBe(false);

    tickFolderPreview();

    // a page opened later in the same browser
    const reopened = TestBed.createComponent(MngFilesComponent);
    reopened.detectChanges();
    await flushAsync();
    reopened.detectChanges();

    expect(isChecked('folderPreview', reopened.nativeElement)).toBe(true);
  });

  it('stays off when the stored value cannot be read', () => {
    spyOn(Storage.prototype, 'getItem').and.throwError('storage is blocked');

    open([dirRow('dir-1', 'Album')]);

    expect(isChecked('folderPreview')).toBe(false);
  });

  it('draws a known cover inside the folder shape when it is on', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));

    tickFolderPreview();

    const cell = thumbnailCell();
    const glyph = glyphAt(0);
    const cover = glyph.querySelector('img');
    expect(glyph).not.toBeNull();
    expect(cover).not.toBeNull();
    expect(cover.getAttribute('src')).toBe(COVER_URL);
    // the folder shape is the only frame: never a glyph on top of a picture
    expect(cell.querySelector('i.bi')).toBeNull();
  });

  it('falls back to the empty folder shape when a cached cover fails to load', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    tickFolderPreview();
    const before = glyphAt(0).getBoundingClientRect();

    glyphAt(0).querySelector('img').dispatchEvent(new Event('error'));
    fixture.detectChanges();

    const cell = thumbnailCell();
    const glyph = glyphAt(0);
    expect(glyph).not.toBeNull();
    expect(glyph.querySelector('img')).toBeNull();
    expect(cell.querySelector('i.bi')).toBeNull();
    const after = glyph.getBoundingClientRect();
    expect(after.width).toBeCloseTo(before.width, 0);
    expect(after.height).toBeCloseTo(before.height, 0);
    // forgotten, so a later listing does not draw the same broken cover again
    expect(component.dirThumbnailCache.has('dir-1')).toBe(false);
  });

  it('keeps the folder shape and adds the cover when the batch answer arrives later', () => {
    open([dirRow('dir-1', 'Album')]);
    tickFolderPreview();
    expect(glyphAt(0)).not.toBeNull();
    expect(glyphAt(0).querySelector('img')).toBeNull();

    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    fixture.detectChanges();

    expect(glyphAt(0).querySelector('img')).not.toBeNull();
  });

  it('draws the empty folder body for a directory whose cover is unknown while it is on', () => {
    open([dirRow('dir-1', 'Album')]);

    tickFolderPreview();

    const cell = thumbnailCell();
    expect(glyphAt(0)).not.toBeNull();
    expect(glyphAt(0).querySelector('img')).toBeNull();
    expect(cell.querySelector('i.bi')).toBeNull();
  });

  it('keeps the folder shape with an empty body when it is turned off', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    tickFolderPreview();
    expect(glyphAt(0).querySelector('img')).not.toBeNull();

    tickFolderPreview();

    expect(glyphAt(0)).not.toBeNull();
    expect(glyphAt(0).querySelector('img')).toBeNull();
    expect(thumbnailCell().querySelector('i.bi')).toBeNull();
  });

  it('draws the smaller body and cover behind a thin folder rail', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));

    tickFolderPreview();

    const glyph = glyphAt(0);
    const box = glyph.getBoundingClientRect();
    const cover = glyph.querySelector('img').getBoundingClientRect();

    // body 54x34 with a 50x30 cover inside the 2px folder-coloured rail
    expect(box.width).toBeCloseTo(54, 0);
    expect(box.height).toBeCloseTo(34, 0);
    expect(cover.width).toBeCloseTo(50, 0);
    expect(cover.height).toBeCloseTo(30, 0);
    expect(cover.left - box.left).toBeCloseTo(2, 0);
    expect(cover.top - box.top).toBeCloseTo(2, 0);
  });

  it('holds a fixed 60px preview slot with the 54x44 silhouette inside it', () => {
    open([dirRow('dir-1', 'Album')]);

    const cell = thumbnailCell();
    const slot = cell.querySelector('.dir-glyph-slot') as HTMLElement;
    const glyph = glyphAt(0);

    expect(slot).not.toBeNull();
    const slotBox = slot.getBoundingClientRect();
    const glyphBox = glyph.getBoundingClientRect();

    // the row rhythm comes from this slot, not from the glyph
    expect(slotBox.height).toBeCloseTo(60, 0);
    expect(glyphBox.width).toBeCloseTo(54, 0);
    expect(glyphBox.height).toBeCloseTo(34, 0);
    // the tab is 10px proud of the body: silhouette 54x44, centred in the slot
    expect(glyphBox.top - 10).toBeGreaterThanOrEqual(slotBox.top - 0.5);
    expect(glyphBox.bottom).toBeLessThanOrEqual(slotBox.bottom + 0.5);
  });

  it('keeps a covered directory the same overall size with the toggle on and off', () => {
    open([dirRow('dir-1', 'Album')]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    const off = glyphAt(0).getBoundingClientRect();
    const offRow = rowAt(0).getBoundingClientRect().height;

    tickFolderPreview();
    const on = glyphAt(0).getBoundingClientRect();
    const onRow = rowAt(0).getBoundingClientRect().height;

    tickFolderPreview();
    const offAgain = glyphAt(0).getBoundingClientRect();

    expect(on.width).toBeCloseTo(off.width, 0);
    expect(on.height).toBeCloseTo(off.height, 0);
    expect(offAgain.width).toBeCloseTo(off.width, 0);
    expect(offAgain.height).toBeCloseTo(off.height, 0);
    expect(onRow).toBeCloseTo(offRow, 0);
  });

  it('gives a covered directory, an uncovered directory and an image file the same row height', async () => {
    open([
      dirRow('dir-1', 'Album'),
      dirRow('dir-2', 'Screenshots'),
      fileRow('f-1', 'photo.jpg', trustedImage(IMAGE_DATA_URL)),
    ]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    tickFolderPreview();
    await settleImages();

    expect(glyphAt(0).querySelector('img')).not.toBeNull();
    expect(glyphAt(1).querySelector('img')).toBeNull();

    atViewportWidth(1200, (viewportWidth) => {
      expect(viewportWidth).toBeGreaterThan(768);

      const heights = [0, 1, 2].map(
        (i) => rowAt(i).getBoundingClientRect().height
      );
      expect(Math.max(...heights) - Math.min(...heights)).toBeLessThanOrEqual(0.5);
      // the tab stays inside the row's own flow box, clear of the row above
      const cell = cellAt(0).getBoundingClientRect();
      const glyphBox = glyphAt(0).getBoundingClientRect();
      expect(glyphBox.top - 10).toBeGreaterThanOrEqual(cell.top - 0.5);
      expect(rowAt(1).getBoundingClientRect().top).toBeGreaterThanOrEqual(
        rowAt(0).getBoundingClientRect().bottom - 0.5
      );
      expect(glyphAt(1).getBoundingClientRect().top - 10).toBeGreaterThanOrEqual(
        rowAt(0).getBoundingClientRect().bottom - 0.5
      );
    });
  });

  it('changes no row height when the preview toggle is switched', async () => {
    open([
      dirRow('dir-1', 'Album'),
      dirRow('dir-2', 'Screenshots'),
      fileRow('f-1', 'photo.jpg', trustedImage(IMAGE_DATA_URL)),
    ]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    await settleImages();

    atViewportWidth(1200, (viewportWidth) => {
      expect(viewportWidth).toBeGreaterThan(768);

      const before = [0, 1, 2].map(
        (i) => rowAt(i).getBoundingClientRect().height
      );
      tickFolderPreview();
      const after = [0, 1, 2].map(
        (i) => rowAt(i).getBoundingClientRect().height
      );

      expect(glyphAt(0).querySelector('img')).not.toBeNull();
      before.forEach((height, i) =>
        expect(Math.abs(after[i] - height)).toBeLessThanOrEqual(0.5)
      );
    });
  });

  it('keeps the same footprint in the compact layout without clipping the shape', async () => {
    mobile$.next(true);
    open([
      dirRow('dir-1', 'Album'),
      dirRow('dir-2', 'Screenshots'),
      fileRow('f-1', 'photo.jpg', trustedImage(IMAGE_DATA_URL)),
    ]);
    component.dirThumbnailCache.set('dir-1', trustedImage(COVER_URL));
    tickFolderPreview();
    await settleImages();

    atPhoneWidth((viewportWidth) => {
      expect(viewportWidth).toBeLessThanOrEqual(768);

      const heights = [0, 1, 2].map(
        (i) => rowAt(i).getBoundingClientRect().height
      );
      expect(Math.max(...heights) - Math.min(...heights)).toBeLessThanOrEqual(0.5);

      const cell = cellAt(0).getBoundingClientRect();
      const glyph = glyphAt(0).getBoundingClientRect();
      expect(glyph.width).toBeCloseTo(54, 0);
      expect(glyph.height).toBeCloseTo(34, 0);
      expect(glyph.left).toBeGreaterThanOrEqual(cell.left - 0.5);
      expect(glyph.right).toBeLessThanOrEqual(cell.right + 0.5);
      // tab included (10px proud): nothing spills above or below the cell
      expect(glyph.top - 10).toBeGreaterThanOrEqual(cell.top - 0.5);
      expect(glyph.bottom).toBeLessThanOrEqual(cell.bottom + 0.5);
    });
  });

  it('leaves file rows exactly as they were under both states', () => {
    open([fileRow('f-1', 'photo.jpg', 'fstore/file/raw?key=file-thumb')]);
    const before = thumbnailCell().innerHTML;

    tickFolderPreview();
    expect(thumbnailCell().innerHTML).toBe(before);

    tickFolderPreview();
    expect(thumbnailCell().innerHTML).toBe(before);
  });

  it('issues no request when it is toggled', () => {
    open([dirRow('dir-1', 'Album')]);
    const posts = http.posts.length;
    const gets = http.gets.length;

    tickFolderPreview();

    expect(http.posts.length).toBe(posts);
    expect(http.gets.length).toBe(gets);
  });

  it('sits beside the tooltip checkbox in the wide layout', () => {
    open([dirRow('dir-1', 'Album')]);

    const labels = Array.from<HTMLElement>(
      fixture.nativeElement.querySelectorAll('mat-checkbox')
    ).map((box) => box.textContent.trim());

    expect(labels).toEqual(['thumbnailTooltip', 'folderPreview']);
  });

  it('keeps the checkbox on screen in the compact layout', () => {
    mobile$.next(true);

    open([dirRow('dir-1', 'Album')]);

    atPhoneWidth((viewportWidth) => {
      expect(viewportWidth).toBeLessThanOrEqual(768);

      const box = checkbox('folderPreview');
      const strip = box.parentElement;
      const boxRect = box.getBoundingClientRect();
      const stripRect = strip.getBoundingClientRect();

      expect(boxRect.width).toBeGreaterThan(0);
      // inside the strip's visible area, not scrolled out to the right
      expect(boxRect.left).toBeGreaterThanOrEqual(stripRect.left);
      expect(boxRect.right).toBeLessThanOrEqual(stripRect.right);
      // and inside the phone viewport
      expect(boxRect.left).toBeGreaterThanOrEqual(0);
      expect(boxRect.right).toBeLessThanOrEqual(viewportWidth);
    });
  });
});
