import { HttpClient } from '@angular/common/http';
import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { UserInfo } from 'src/common/user-info';
import { UserService } from './user.service';

@Injectable({
  providedIn: 'root',
})
export class WebSocketNotificationService implements OnDestroy {
  private ws: WebSocket | null = null;
  private countSubject = new BehaviorSubject<number>(0);
  private reconnectAttempts = 0;
  private maxReconnectDelay = 30000; // 30 seconds
  private baseReconnectDelay = 1000; // 1 second
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private ticketSubscription: Subscription | null = null;
  private userInfoSubscription: Subscription;
  private destroyed = false;
  private sessionActive = false;

  /** Observable of current notification count */
  count$: Observable<number> = this.countSubject.asObservable();

  constructor(private http: HttpClient, private userService: UserService) {
    this.userInfoSubscription = this.userService.userInfoObservable.subscribe(
      (userInfo: UserInfo) => {
        if (userInfo) {
          this.startSession();
        } else {
          this.endSession();
        }
      }
    );
  }

  /** Get current count value */
  get currentCount(): number {
    return this.countSubject.value;
  }

  private get canSubscribe(): boolean {
    return !this.destroyed && this.sessionActive;
  }

  private startSession(): void {
    if (this.destroyed || this.sessionActive) {
      return;
    }

    this.sessionActive = true;
    this.reconnectAttempts = 0;
    this.connect();
  }

  private endSession(): void {
    const hadSession = this.sessionActive;
    this.sessionActive = false;
    this.cancelReconnect();
    this.clearTicketRequest();

    if (this.ws) {
      this.ws.close(1000, 'session ended');
      this.ws = null;
    }

    this.reconnectAttempts = 0;
    if (hadSession) {
      this.countSubject.next(0);
    }
  }

  private connect(): void {
    if (!this.canSubscribe) {
      return;
    }

    // Fetch a one-time websocket ticket from user-vault
    this.ticketSubscription = this.http
      .post<any>(`user-vault/open/api/v1/notification/ws-ticket`, {})
      .subscribe({
        next: (resp) => {
          this.ticketSubscription = null;
          if (!this.canSubscribe) {
            return;
          }

          if (resp.error || !resp.data?.ticket) {
            console.error('[WS-Notification] Failed to get ticket', resp);
            this.scheduleReconnect();
            return;
          }

          this.openSocket(resp.data.ticket);
        },
        error: (err) => {
          this.ticketSubscription = null;
          if (!this.canSubscribe) {
            return;
          }

          console.error('[WS-Notification] Failed to get ticket', err);
          this.scheduleReconnect();
        },
      });
  }

  private openSocket(ticket: string): void {
    // Build WebSocket URL with the ticket
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const wsUrl = `${protocol}//${host}/user-vault/open/api/v2/notification/ws?token=${encodeURIComponent(
      ticket
    )}`;

    console.log('[WS-Notification] Connecting...');
    const socket = new WebSocket(wsUrl);
    this.ws = socket;

    socket.onopen = () => {
      console.log('[WS-Notification] Connected');
      this.reconnectAttempts = 0;
    };

    socket.onmessage = (event: MessageEvent) => {
      if (socket !== this.ws) {
        return;
      }
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'count' && typeof data.data === 'number') {
          this.countSubject.next(data.data);
        }
      } catch (e) {
        console.error('[WS-Notification] Failed to parse message', e);
      }
    };

    socket.onclose = (event: CloseEvent) => {
      console.log(
        `[WS-Notification] Closed (code=${event.code}, clean=${event.wasClean})`
      );
      if (socket !== this.ws) {
        return;
      }
      this.ws = null;
      if (this.canSubscribe) {
        this.scheduleReconnect();
      }
    };

    socket.onerror = () => {
      // onclose will be called after onerror
      console.error('[WS-Notification] Error');
    };
  }

  private scheduleReconnect(): void {
    if (!this.canSubscribe) {
      return;
    }

    const delay = Math.min(
      this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts),
      this.maxReconnectDelay
    );
    this.reconnectAttempts++;

    console.log(
      `[WS-Notification] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`
    );

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  private cancelReconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }

  private clearTicketRequest(): void {
    if (this.ticketSubscription) {
      this.ticketSubscription.unsubscribe();
      this.ticketSubscription = null;
    }
  }

  ngOnDestroy(): void {
    this.destroyed = true;
    this.userInfoSubscription.unsubscribe();
    this.cancelReconnect();
    this.clearTicketRequest();

    if (this.ws) {
      this.ws.close(1000, 'component destroyed');
      this.ws = null;
    }
  }
}
