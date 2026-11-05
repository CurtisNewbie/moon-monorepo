import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { distinctUntilChanged } from 'rxjs/operators';

/**
 * Compact ends at this width, inclusive, and wide starts at the next pixel.
 * src/styles.css must use the same two widths; scripts/check-breakpoints.js
 * fails when this constant and the stylesheet blocks disagree.
 */
export const MOBILE_MAX_WIDTH = 768;

/** The single width query the whole app asks. */
export const MOBILE_QUERY = `(max-width: ${MOBILE_MAX_WIDTH}px)`;

/** Old Safari's MediaQueryList, before addEventListener existed. */
interface LegacyMediaQueryList {
  addListener(listener: (mediaQueryList: MediaQueryList) => void): void;
  removeListener(listener: (mediaQueryList: MediaQueryList) => void): void;
}

@Injectable({
  providedIn: 'root',
})
export class Env implements OnDestroy {
  private readonly mediaQuery: MediaQueryList = window.matchMedia(MOBILE_QUERY);
  private readonly currentMobile = new BehaviorSubject<boolean>(this.mediaQuery.matches);

  /**
   * Emits the current answer immediately, then only when the answer changes.
   * One singleton listener feeds this however many subscribers there are.
   */
  readonly mobile$: Observable<boolean> = this.currentMobile.pipe(distinctUntilChanged());

  private readonly onChange = (event: MediaQueryListEvent) => this.currentMobile.next(event.matches);
  private readonly onLegacyChange = () => this.currentMobile.next(this.mediaQuery.matches);

  constructor() {
    if (typeof this.mediaQuery.addEventListener === 'function') {
      this.mediaQuery.addEventListener('change', this.onChange);
    } else {
      (this.mediaQuery as unknown as LegacyMediaQueryList).addListener(this.onLegacyChange);
    }
  }

  /**
   * Whether the compact layout owns the viewport, read from the media query
   * rather than a cached field or window.innerWidth.
   */
  isMobile(): boolean {
    return this.mediaQuery.matches;
  }

  ngOnDestroy(): void {
    if (typeof this.mediaQuery.removeEventListener === 'function') {
      this.mediaQuery.removeEventListener('change', this.onChange);
    } else {
      (this.mediaQuery as unknown as LegacyMediaQueryList).removeListener(this.onLegacyChange);
    }
  }
}
