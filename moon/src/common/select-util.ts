export interface Option<T> {
  name: string;
  value: T;
}

/** filter candidates that contains the value */
export function filterAlike(candidates: string[], value: string): string[] {
  if (!value) return candidates;

  return candidates.filter((option) =>
    option.toLowerCase().includes(value.toLowerCase())
  );
}
