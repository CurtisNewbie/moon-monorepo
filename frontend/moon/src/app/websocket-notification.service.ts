import { HttpClient } from "@angular/common/http";
import { Injectable, OnDestroy } from "@angular/core";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable({
  providedIn: "root",
})
export class WebSocketNotificationService implements OnDestroy {
  private ws: WebSocket | null = null;
  private countSubject = new BehaviorSubject<number>(0);
  private reconnectAttempts = 0;
  private maxReconnectDelay = 30000; // 30 seconds
  private baseReconnectDelay = 1000; // 1 second
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private destroyed = false;
  private intentionalClose = false;

  /** Observable of current notification count */
  count$: Observable<number> = this.countSubject.asObservable();

  constructor(private http: HttpClient) {
    this.connect();
  }

  /** Get current count value */
  get currentCount(): number {
    return this.countSubject.value;
  }

  private connect(): void {
    if (this.destroyed) return;

    // Fetch a one-time websocket ticket from user-vault
    this.http.post<any>(`user-vault/open/api/v1/notification/ws-ticket`, {}).subscribe({
      next: (resp) => {
        if (resp.error || !resp.data?.ticket) {
          console.error("[WS-Notification] Failed to get ticket", resp);
          this.scheduleReconnect();
          return;
        }

        const ticket = resp.data.ticket;

        // Build WebSocket URL with the ticket
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const host = window.location.host;
        const wsUrl = `${protocol}//${host}/user-vault/open/api/v2/notification/ws?token=${encodeURIComponent(ticket)}`;

        console.log("[WS-Notification] Connecting...");
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log("[WS-Notification] Connected");
          this.reconnectAttempts = 0;
          this.intentionalClose = false;
        };

        this.ws.onmessage = (event: MessageEvent) => {
          try {
            const data = JSON.parse(event.data);
            if (data.type === "count" && typeof data.data === "number") {
              this.countSubject.next(data.data);
            }
          } catch (e) {
            console.error("[WS-Notification] Failed to parse message", e);
          }
        };

        this.ws.onclose = (event: CloseEvent) => {
          console.log(`[WS-Notification] Closed (code=${event.code}, clean=${event.wasClean})`);
          this.ws = null;
          if (!this.intentionalClose && !this.destroyed) {
            this.scheduleReconnect();
          }
        };

        this.ws.onerror = () => {
          // onclose will be called after onerror
          console.error("[WS-Notification] Error");
        };
      },
      error: (err) => {
        console.error("[WS-Notification] Failed to get ticket", err);
        this.scheduleReconnect();
      },
    });
  }

  private scheduleReconnect(): void {
    if (this.destroyed) return;

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

  ngOnDestroy(): void {
    this.destroyed = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.intentionalClose = true;
      this.ws.close(1000, "component destroyed");
      this.ws = null;
    }
  }
}
