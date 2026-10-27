import { FakeMediaQueryList, installFakeMatchMedia, useFakeMatchMedia } from 'src/testing/fake-media-query-list';
import { Env } from './env-util';

describe('Env', () => {
  let fake: FakeMediaQueryList;

  function createEnv(matches: boolean): Env {
    fake = installFakeMatchMedia(matches);
    return new Env();
  }

  function record(env: Env): boolean[] {
    const values: boolean[] = [];
    env.mobile$.subscribe((value) => values.push(value));
    return values;
  }

  it('answers isMobile from the match the media query already had', () => {
    const env = createEnv(true);

    expect(env.isMobile()).toBeTrue();
    expect(fake.media).toBe('(max-width: 768px)');
  });

  it('starts mobile$ with the current value', () => {
    const env = createEnv(true);

    expect(record(env)).toEqual([true]);
  });

  it('reports a real change from both isMobile and mobile$', () => {
    const env = createEnv(true);
    const values = record(env);

    fake.fireChange(false);

    expect(env.isMobile()).toBeFalse();
    expect(values).toEqual([true, false]);
  });

  it('stays silent when a change event reports the value it already had', () => {
    const env = createEnv(false);
    const values = record(env);

    fake.fireChange(false);

    expect(env.isMobile()).toBeFalse();
    expect(values).toEqual([false]);
  });

  it('adds one listener however many subscribers mobile$ has', () => {
    const env = createEnv(true);

    record(env);
    record(env);

    expect(fake.listenerCount).toBe(1);
  });

  it('falls back to addListener when the list has no addEventListener', () => {
    fake = useFakeMatchMedia(FakeMediaQueryList.withoutAddEventListener(true));
    const env = new Env();
    const values = record(env);

    fake.fireChange(false);

    expect(fake.legacyListenerCount).toBe(1);
    expect(values).toEqual([true, false]);
  });
});
