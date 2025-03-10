import { Component, Inject, OnInit } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { environment } from "src/environments/environment";
import { WPath } from "../manage-paths/manage-paths.component";
import { MngResDialogComponent } from "../mng-res-dialog/mng-res-dialog.component";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface DialogDat {
  path: WPath;
}

@Component({
  selector: "app-mng-path-dialog",
  templateUrl: "./mng-path-dialog.component.html",
  styleUrls: ["./mng-path-dialog.component.css"],
})
export class MngPathDialogComponent implements OnInit {
  bindToResCode = "";
  PATH_TYPES = [
    { val: "PROTECTED", name: "Protected" },
    { val: "PUBLIC", name: "Public" },
  ];

  constructor(
    public dialogRef: MatDialogRef<MngResDialogComponent, DialogDat>,
    @Inject(MAT_DIALOG_DATA) public dat: DialogDat,
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {}

  update() {
    this.http
      .post<any>(`user-vault/open/api/path/update`, {
        type: this.dat.path.ptype,
        pathNo: this.dat.path.pathNo,
        group: this.dat.path.pgroup,
        resCode: this.bindToResCode,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
        },
        complete: () => {
          this.dialogRef.close();
        },
      });
  }
}
