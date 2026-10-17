import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { HttpClient } from '@angular/common/http';
import { fakeAsync, TestBed, tick } from '@angular/core/testing';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { UserInfo } from 'src/common/user-info';

import { UserService } from './user.service';
import { WebSocketNotificationService } from './websocket-notification.service';

const WS_TICKET_URL = 'user-vault/open/api/v1/notification/ws-ticket';

class FakeWebSocket {
  static instances: FakeWebSocket[] = [];

  url: string;
  onopen: ((event: any) => void) | null = null;
  onmessage: ((event: any) => void) | null = null;
  onclose: ((event: any) => void) | null = null;
  onerror: ((event: any) => void) | null = null;
  closeCalls: Array<{ code?: number; reason?: string }> = [];

  constructor(url: string) {
    this.url = url;
    FakeWebSocket.instances.push(this);
  }

  close(code?: number, reason?: string): void {
    this.closeCalls.push({ code, reason });
  }

  triggerMessage(payload: any): void {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(payload) } as MessageEvent);
    }
  }

  triggerClose(code = 1006, wasClean = false): void {
    if (this.onclose) {
      this.onclose({ code, wasClean } as CloseEvent);
    }
  }
}

class UserServiceStub {
  private userInfoSubject = new BehaviorSubject<UserInfo>(null);
  userInfoObservable: Observable<UserInfo> =
    this.userInfoSubject.asObservable();

  emit(userInfo: UserInfo): void {
    this.userInfoSubject.next(userInfo);
  }
}

class LazyUserServiceStub {
  private userInfoSubject = new Subject<UserInfo>();
  userInfoObservable: Observable<UserInfo> =
    this.userInfoSubject.asObservable();

  emit(userInfo: UserInfo): void {
    this.userInfoSubject.next(userInfo);
  }
}

describe('WebSocketNotificationService Login Session gate', () => {
  let service: WebSocketNotificationService;
  let httpMock: HttpTestingController;
  let userService: UserServiceStub;
  const realWebSocket = window.WebSocket;

  beforeEach(() => {
    FakeWebSocket.instances = [];
    userService = new UserServiceStub();
    (window as any).WebSocket = FakeWebSocket;

    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [{ provide: UserService, useValue: userService }],
    });
    service = TestBed.inject(WebSocketNotificationService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    service.ngOnDestroy();
    (window as any).WebSocket = realWebSocket;
    httpMock.verify();
  });

  it('makes no ticket request before identity is confirmed', fakeAsync(() => {
    tick(60000);

    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(0);
  }));

  it('does not publish a new count for a null identity outside a session', fakeAsync(() => {
    const lazyUserService = new LazyUserServiceStub();
    const lazyService = new WebSocketNotificationService(
      TestBed.inject(HttpClient),
      lazyUserService as any
    );
    const counts: number[] = [];
    lazyService.count$.subscribe((count) => counts.push(count));

    lazyUserService.emit(null);

    expect(counts).toEqual([0]);
    lazyService.ngOnDestroy();
  }));

  it('requests one ticket and opens one socket carrying it after identity is confirmed', fakeAsync(() => {
    httpMock.expectNone(WS_TICKET_URL);

    userService.emit({ id: 1, username: 'u' } as UserInfo);
    userService.emit({ id: 1, username: 'u' } as UserInfo);

    const req = httpMock.expectOne(WS_TICKET_URL);
    expect(req.request.method).toBe('POST');
    req.flush({ error: false, data: { ticket: 'ticket-1' } });

    expect(FakeWebSocket.instances.length).toBe(1);
    const socket = FakeWebSocket.instances[0];
    expect(socket.url).toContain(
      `/user-vault/open/api/v2/notification/ws?token=${encodeURIComponent(
        'ticket-1'
      )}`
    );

    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(1);
  }));

  it('closes the socket and publishes zero when identity becomes null', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-1' } });
    const socket = FakeWebSocket.instances[0];

    const counts: number[] = [];
    service.count$.subscribe((count) => counts.push(count));
    socket.triggerMessage({ type: 'count', data: 7 });
    expect(counts[counts.length - 1]).toBe(7);

    userService.emit(null);

    expect(socket.closeCalls.length).toBe(1);
    expect(counts[counts.length - 1]).toBe(0);

    socket.triggerMessage({ type: 'count', data: 9 });
    expect(counts[counts.length - 1]).toBe(0);

    socket.triggerClose(1000, true);
    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(1);
  }));

  it('retries a transport failure after the backoff delay', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .error(new ErrorEvent('network error'));

    tick(999);
    httpMock.expectNone(WS_TICKET_URL);

    tick(1);
    const retry = httpMock.expectOne(WS_TICKET_URL);
    retry.flush({ error: false, data: { ticket: 'ticket-2' } });
    expect(FakeWebSocket.instances.length).toBe(1);
  }));

  it('retries a server error after the backoff delay', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ msg: 'boom' }, { status: 500, statusText: 'Server Error' });

    tick(999);
    httpMock.expectNone(WS_TICKET_URL);

    tick(1);
    const retries = httpMock.match(WS_TICKET_URL);
    expect(retries.length).toBe(1);
    retries[0].flush({ error: false, data: { ticket: 'ticket-2' } });
    expect(FakeWebSocket.instances.length).toBe(1);
  }));

  it('makes no further ticket request after a 401 rejection', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush(
        { msg: 'unauthorized' },
        { status: 401, statusText: 'Unauthorized' }
      );

    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(0);
  }));

  it('makes no further ticket request after a 403 rejection', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush(
        { msg: 'forbidden' },
        { status: 403, statusText: 'Forbidden' }
      );

    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(0);
  }));

  it('connects again after a terminal rejection when a new session starts', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush(
        { msg: 'unauthorized' },
        { status: 401, statusText: 'Unauthorized' }
      );
    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);

    userService.emit(null);
    userService.emit({ id: 1, username: 'u' } as UserInfo);

    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-2' } });

    expect(FakeWebSocket.instances.length).toBe(1);
    expect(FakeWebSocket.instances[0].url).toContain(
      `token=${encodeURIComponent('ticket-2')}`
    );

    tick(60000);
    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(1);
  }));

  it('reconnects after a transport-failure socket close inside the session', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-1' } });
    const socket = FakeWebSocket.instances[0];

    socket.triggerClose(1006, false);

    tick(999);
    httpMock.expectNone(WS_TICKET_URL);

    tick(1);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-2' } });

    expect(FakeWebSocket.instances.length).toBe(2);
    expect(FakeWebSocket.instances[1].url).toContain(
      `token=${encodeURIComponent('ticket-2')}`
    );
  }));

  it('discards a retry pending when the session ends', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .error(new ErrorEvent('network error'));
    tick(999);
    httpMock.expectNone(WS_TICKET_URL);

    userService.emit(null);
    tick(60000);

    httpMock.expectNone(WS_TICKET_URL);
    expect(FakeWebSocket.instances.length).toBe(0);
  }));

  it('connects again after logging back in', fakeAsync(() => {
    userService.emit({ id: 1, username: 'u' } as UserInfo);
    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-1' } });
    expect(FakeWebSocket.instances.length).toBe(1);

    userService.emit(null);
    userService.emit({ id: 1, username: 'u' } as UserInfo);

    httpMock
      .expectOne(WS_TICKET_URL)
      .flush({ error: false, data: { ticket: 'ticket-2' } });

    expect(FakeWebSocket.instances.length).toBe(2);
    expect(FakeWebSocket.instances[1].url).toContain(
      `token=${encodeURIComponent('ticket-2')}`
    );
  }));
});
