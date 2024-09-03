import { Injectable } from '@angular/core';
import { Observable, Subject, Subscription, timer } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class PlatformNotificationService {

  changeSubject = new Subject();
  private timerSubscription: Subscription;

  constructor() {
    this.timerSubscription = timer(1000, 3000).subscribe(event => {
      this.triggerChange();
    });
  }

  subscribeChange(): Observable<any> {
    return this.changeSubject.asObservable();
  }

  triggerChange() {
    this.changeSubject.next();
  }

}
