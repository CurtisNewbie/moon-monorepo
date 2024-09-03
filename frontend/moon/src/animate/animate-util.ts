import {
  animate,
  AnimationTriggerMetadata,
  state,
  style,
  transition,
  trigger,
} from "@angular/animations";

export function animateElementExpanding(): AnimationTriggerMetadata {
  return trigger("detailExpand", [
    state("collapsed", style({ height: "0px", minHeight: "0" })),
    state("expanded", style({ height: "*" })),
    transition(
      "expanded <=> collapsed",
      animate("100ms cubic-bezier(0.4, 0.0, 0.2, 1)")
    ),
  ]);
}

export function copy<T>(t: T): T {
  if (t == null) return null;
  return { ...t }
}

export function isIdEqual(t, v): boolean {
  if (t == null || v == null) return false;
  return t.id === v.id;
}

export function getExpanded(row, currExpanded, closeForMobile: boolean = false) {
  if (closeForMobile) return null;
  // null means row is the expanded one, so we return null to make it collapsed
  return isIdEqual(row, currExpanded) ? null : copy(row);
}
