import { Injectable, OnDestroy } from "@angular/core";
import { Observable, Subject, Subscription } from "rxjs";
import { WebSocketNotificationService } from "./websocket-notification.service";

@Injectable({
  providedIn: "root",
})
export class PlatformNotificationService implements OnDestroy {
  private changeSubject = new Subject<void>();
  private wsSubscription: Subscription;

  constructor(private wsService: WebSocketNotificationService) {
    // When WS pushes a new count, notify subscribers
    this.wsSubscription = this.wsService.count$.subscribe(() => {
      this.changeSubject.next();
    });
  }

  /** Subscribe to notification change events */
  subscribeChange(): Observable<void> {
    return this.changeSubject.asObservable();
  }

  /** Trigger a manual refresh (e.g., after marking notifications as read) */
  triggerChange(): void {
    this.changeSubject.next();
  }

  /** Get current unread count (from WebSocket) */
  get currentCount(): number {
    return this.wsService.currentCount;
  }

  ngOnDestroy(): void {
    this.wsSubscription?.unsubscribe();
  }
}
