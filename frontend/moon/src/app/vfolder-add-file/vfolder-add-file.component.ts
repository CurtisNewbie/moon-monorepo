import { Component, Inject, OnInit } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { ConfirmDialog } from "src/common/dialog";
import { VFolderBrief } from "src/common/folder";
import { filterAlike } from "src/common/select-util";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

type VfFile = {
  fileKey: string;
  name: string;
};

type Data = {
  files: VfFile[];
};

@Component({
  selector: "app-vfolder-add-file",
  templateUrl: "./vfolder-add-file.component.html",
  styleUrls: ["./vfolder-add-file.component.css"],
})
export class VfolderAddFileComponent implements OnInit {
  addToVFolderName: string;

  /** list of brief info of all vfolder that we created */
  vfolderBrief: VFolderBrief[] = [];
  /** Auto complete for vfolders that we may add file into */
  autoCompAddToVFolderName: string[];

  constructor(
    public dialogRef: MatDialogRef<VfolderAddFileComponent, Data>,
    @Inject(MAT_DIALOG_DATA) public dat: Data,
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.fetchOwnedVFolderBrief();
  }

  onAddToVFolderNameChanged() {
    this.autoCompAddToVFolderName = filterAlike(
      this.vfolderBrief.map((v) => v.name),
      this.addToVFolderName
    );
  }

  fetchOwnedVFolderBrief() {
    this.http.get<any>(`vfm/open/api/vfolder/brief/owned`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.vfolderBrief = resp.data;
        this.onAddToVFolderNameChanged();
      },
    });
  }

  addToVirtualFolder() {
    const vfolderName = this.addToVFolderName;
    if (!vfolderName) {
      this.snackBar.open("Please select a folder first", "ok", {
        duration: 3000,
      });
      return;
    }

    let addToFolderNo;
    let matched: VFolderBrief[] = this.vfolderBrief.filter(
      (v) => v.name === vfolderName
    );
    if (!matched || matched.length < 1) {
      this.snackBar.open(
        "Virtual Folder not found, please check and try again",
        "ok",
        { duration: 3000 }
      );
      return;
    }
    if (matched.length > 1) {
      this.snackBar.open(
        "Found multiple virtual folder with the same name, please try again",
        "ok",
        { duration: 3000 }
      );
      return;
    }
    addToFolderNo = matched[0].folderNo;

    this.http
      .post<any>(`vfm/open/api/vfolder/file/add`, {
        folderNo: addToFolderNo,
        fileKeys: this.dat.files.map((f) => f.fileKey),
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.snackBar.open("Success", "ok", { duration: 3000 });
        },
      });
  }
}
