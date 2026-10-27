import { Injectable } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { ConfirmDialogComponent } from "src/app/dialog/confirm/confirm-dialog.component";
import { dialogSize } from "src/common/dialog-size";

@Injectable({
  providedIn: 'root'
})
export class ConfirmDialog {

  constructor(private dialog: MatDialog) { }

  show(title: string, msg: string[], onConfirm, widthPx: number = 500) {
    this.dialog.open(ConfirmDialogComponent, {
      ...dialogSize(widthPx),
      data: {
        title: title,
        msg: msg,
        isNoBtnDisplayed: true,
      },
    })
      .afterClosed()
      .subscribe((confirm) => {
        if (!confirm) {
          return;
        }
        onConfirm()
      });
  }

}