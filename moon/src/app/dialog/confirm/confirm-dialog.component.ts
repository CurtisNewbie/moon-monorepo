import { Component, Inject } from "@angular/core";
import { MatDialogRef } from "@angular/material/dialog";
import { MAT_DIALOG_DATA } from "@angular/material/dialog";

export interface ConfirmDialogData {
  title: string;
  msg: string[];
  isNoBtnDisplayed: boolean;
}

@Component({
  selector: "confirm-dialog-component",
  template: `
    <h1 mat-dialog-title>{{data.title}}
    </h1>
    <div mat-dialog-content>
        <p *ngFor="let line of data.msg">{{ line }}</p>
    </div>
    <div mat-dialog-actions>
        <button mat-button [mat-dialog-close]="true">Yes</button>
        <button mat-button *ngIf="data.isNoBtnDisplayed" [mat-dialog-close]="false" cdkFocusInitial>No</button>
    </div>
  `,
})
export class ConfirmDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<ConfirmDialogComponent, ConfirmDialogData>,
    @Inject(MAT_DIALOG_DATA) public data: ConfirmDialogData
  ) { }
}
