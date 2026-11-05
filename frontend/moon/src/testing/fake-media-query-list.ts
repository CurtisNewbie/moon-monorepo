type ChangeListener = (event: MediaQueryListEvent) => void;
type LegacyListener = (mediaQueryList: MediaQueryList) => void;

/**
 * Stand-in for the MediaQueryList that window.matchMedia returns. The test owns
 * `matches` and decides when a change happens: `fireChange` sets the value and
 * notifies the listeners the way the browser does, while assigning `matches`
 * alone changes the value with no event.
 *
 * Everything width-dependent is asserted through this one seam.
 */
export class FakeMediaQueryList {
  matches: boolean;
  readonly media: string;
  onchange: ChangeListener | null = null;

  private readonly changeListeners: ChangeListener[] = [];
  private readonly legacyListeners: LegacyListener[] = [];

  constructor(matches: boolean, media: string = '(max-width: 768px)') {
    this.matches = matches;
    this.media = media;
  }

  /**
   * A list shaped like old Safari's: addListener exists, addEventListener does
   * not.
   */
  static withoutAddEventListener(matches: boolean): FakeMediaQueryList {
    const fake = new FakeMediaQueryList(matches);
    (fake as any).addEventListener = undefined;
    return fake;
  }

  /** Listeners registered through addEventListener. */
  get listenerCount(): number {
    return this.changeListeners.length;
  }

  /** Listeners registered through addListener. */
  get legacyListenerCount(): number {
    return this.legacyListeners.length;
  }

  addEventListener(type: string, listener: ChangeListener): void {
    if (type !== 'change') {
      return;
    }
    if (this.changeListeners.indexOf(listener) < 0) {
      this.changeListeners.push(listener);
    }
  }

  removeEventListener(type: string, listener: ChangeListener): void {
    if (type !== 'change') {
      return;
    }
    const index = this.changeListeners.indexOf(listener);
    if (index >= 0) {
      this.changeListeners.splice(index, 1);
    }
  }

  addListener(listener: LegacyListener): void {
    if (this.legacyListeners.indexOf(listener) < 0) {
      this.legacyListeners.push(listener);
    }
  }

  removeListener(listener: LegacyListener): void {
    const index = this.legacyListeners.indexOf(listener);
    if (index >= 0) {
      this.legacyListeners.splice(index, 1);
    }
  }

  /** Sets matches and notifies every listener, as a browser change does. */
  fireChange(matches: boolean = this.matches): void {
    this.matches = matches;
    const event = { matches, media: this.media } as MediaQueryListEvent;

    if (this.onchange) {
      this.onchange(event);
    }
    this.changeListeners.slice().forEach((listener) => listener(event));
    this.legacyListeners.slice().forEach((listener) =>
      listener(this as unknown as MediaQueryList));
  }
}

/** Installs a fake with the given current match on window.matchMedia. */
export function installFakeMatchMedia(matches: boolean): FakeMediaQueryList {
  return useFakeMatchMedia(new FakeMediaQueryList(matches));
}

/** Installs an already built fake on window.matchMedia. */
export function useFakeMatchMedia(fake: FakeMediaQueryList): FakeMediaQueryList {
  spyOn(window, 'matchMedia').and.returnValue(fake as unknown as MediaQueryList);
  return fake;
}
