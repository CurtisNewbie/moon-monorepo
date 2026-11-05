import { MatDialogConfig } from '@angular/material/dialog';

/** The widest a dialog may get, relative to the viewport. */
export const DIALOG_MAX_WIDTH = '95vw';

/** The width and maximum width to spread into a MatDialogConfig. */
export type DialogSize = Pick<MatDialogConfig, 'width' | 'maxWidth'>;

/**
 * Sizes a dialog so it is never wider than its designed pixel width or 95% of
 * the viewport, whichever is smaller. Material's own 80vw default is what clips
 * content on a phone, so maxWidth is raised to the same cap.
 *
 *   this.dialog.open(SomeDialog, { ...dialogSize(800), data });
 */
export function dialogSize(widthPx: number): DialogSize {
  return {
    width: `min(${widthPx}px, ${DIALOG_MAX_WIDTH})`,
    maxWidth: DIALOG_MAX_WIDTH,
  };
}

/**
 * Sizes a dialog that is designed to fill part of the viewport rather than a
 * fixed pixel width: `percent` is the share of the viewport it asks for, and the
 * maximum width is the same cap `dialogSize` uses, so it never reaches the edges.
 *
 *   this.dialog.open(SomeDialog, { ...dialogSizeVw(80), data });
 */
export function dialogSizeVw(percent: number): DialogSize {
  return {
    width: `${percent}vw`,
    maxWidth: DIALOG_MAX_WIDTH,
  };
}
