import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { fakeAsync, TestBed, tick } from '@angular/core/testing';
import { MatSnackBar } from '@angular/material/snack-bar';

import { BrowseHistoryRecorder } from './browse-history';

// Karma installs window.onbeforeunload to detect page reloads; a synthetic
// dispatch would trip it, so suspend it around the dispatch only.
function dispatchBeforeUnload() {
  const karmaUnloadHandler = window.onbeforeunload;
  window.onbeforeunload = null;
  try {
    window.dispatchEvent(new Event('beforeunload'));
  } finally {
    window.onbeforeunload = karmaUnloadHandler;
  }
}

describe('BrowseHistoryRecorder dir-page writes', () => {
  let service: BrowseHistoryRecorder;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [{ provide: MatSnackBar, useValue: { open: () => undefined } }],
    });
    service = TestBed.inject(BrowseHistoryRecorder);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    // TestBed does not destroy providedIn:'root' services, so the beforeunload
    // listener would otherwise leak between specs.
    service.ngOnDestroy();
    httpMock.verify();
  });

  it('serializes dir-page writes and never lands an earlier fileKey', fakeAsync(() => {
    service.recordDirPage('dir', 'f1');
    service.recordDirPage('dir', 'f2');
    service.recordDirPage('dir', 'f3');

    // Nothing may be dispatched while the writes are still being coalesced.
    expect(httpMock.match('/vfm/history/dir/last-page').length).toBe(0);

    tick(300);

    // Exactly one request, and it must carry the latest fileKey.
    const reqs = httpMock.match('/vfm/history/dir/last-page');
    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body).toEqual({ dirKey: 'dir', fileKey: 'f3' });

    reqs[0].flush({});
    tick(300);
  }));

  it('flushes the pending position via sendBeacon on beforeunload', async () => {
    const beacon = spyOn(navigator, 'sendBeacon').and.returnValue(true);

    service.recordDirPage('dir', 'f9');
    dispatchBeforeUnload();

    expect(beacon).toHaveBeenCalledTimes(1);
    const args = beacon.calls.mostRecent().args;
    expect(args[0]).toBe('/vfm/history/dir/last-page');

    const payload = args[1] as Blob;
    expect(payload.type).toBe('application/json');
    expect(JSON.parse(await payload.text())).toEqual({ dirKey: 'dir', fileKey: 'f9' });

    // The beacon path is used instead of the regular HTTP request.
    expect(httpMock.match('/vfm/history/dir/last-page').length).toBe(0);

    // Let the coalescing timer fire, then drain it so it does not linger past the test.
    await new Promise((resolve) => setTimeout(resolve, 350));
    httpMock.match('/vfm/history/dir/last-page').forEach((r) => r.flush({}));
  });

  it('re-beacons the latest position when the tab closes while a write is in flight', async () => {
    const beacon = spyOn(navigator, 'sendBeacon').and.returnValue(true);

    service.recordDirPage('dir', 'f5');

    // Let the debounce fire: the POST is now dispatched and still in flight.
    await new Promise((resolve) => setTimeout(resolve, 350));
    const inFlight = httpMock.match('/vfm/history/dir/last-page');
    expect(inFlight.length).toBe(1);
    expect(beacon).not.toHaveBeenCalled();

    dispatchBeforeUnload();

    expect(beacon).toHaveBeenCalledTimes(1);
    const args = beacon.calls.mostRecent().args;
    expect(args[0]).toBe('/vfm/history/dir/last-page');
    const payload = args[1] as Blob;
    expect(JSON.parse(await payload.text())).toEqual({ dirKey: 'dir', fileKey: 'f5' });

    // Drain the in-flight request that unload would have aborted.
    inFlight.forEach((r) => r.flush({}));
  });
});
