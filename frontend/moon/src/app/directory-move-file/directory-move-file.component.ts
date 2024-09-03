import { Component, Inject, OnInit } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { DirBrief } from "src/common/file-info";
import { filterAlike } from "src/common/select-util";
import { Toaster } from "../notification.service";
import { environment } from "src/environments/environment";
import { HttpClient } from "@angular/common/http";

type DfFile = {
  fileKey: string;
  name: string;
};

type Data = {
  files: DfFile[];
};

@Component({
  selector: "app-directory-move-file",
  templateUrl: "./directory-move-file.component.html",
  styleUrls: ["./directory-move-file.component.css"],
})
export class DirectoryMoveFileComponent implements OnInit {
  /** list of brief info of all directories that we can access */
  dirBriefList: DirBrief[] = [];
  /** auto complete for dirs that we may move file into */
  autoCompMoveIntoDirs: string[] = [];
  /** name of dir that we may move file into */
  moveIntoDirName: string = null;

  onMoveIntoDirNameChanged = () =>
    (this.autoCompMoveIntoDirs = filterAlike(
      this.dirBriefList.map((v) => v.name),
      this.moveIntoDirName
    ));

  constructor(
    public dialogRef: MatDialogRef<DirectoryMoveFileComponent, Data>,
    @Inject(MAT_DIALOG_DATA) public dat: Data,
    private http: HttpClient,
    private toaster: Toaster
  ) {}

  ngOnInit(): void {
    this._fetchDirBriefList();
  }

  // fetch dir brief list
  private _fetchDirBriefList() {
    this.http.get<any>(`${environment.vfm}/open/api/file/dir/list`).subscribe({
      next: (resp) => {
        this.dirBriefList = resp.data;
        this.onMoveIntoDirNameChanged();
      },
    });
  }

  findMoveIntoDirFileKey(dirName: string) {
    let matched: DirBrief[] = this.dirBriefList.filter(
      (v) => v.name === dirName
    );
    if (!matched || matched.length < 1) {
      this.toaster.toast("Directory not found, please check and try again");
      return;
    }
    if (matched.length > 1) {
      this.toaster.toast(
        "Found multiple directories with the same name, please update their names and try again",
        4000
      );
      return;
    }
    return matched[0].uuid;
  }

  moveToDir() {
    const moveIntoDirName = this.moveIntoDirName;
    if (!moveIntoDirName) {
      this.toaster.toast("Please select directory");
      return;
    }
    const key = this.findMoveIntoDirFileKey(moveIntoDirName);
    if (!key) {
      return;
    }
    this._moveEachToDir(this.dat.files, key, 0);
  }

  private _moveEachToDir(files, dirFileKey: string, offset: number) {
    if (offset >= files.length) {
      return;
    }

    let curr = files[offset];
    this.http
      .post(`${environment.vfm}/open/api/file/move-to-dir`, {
        uuid: curr.fileKey,
        parentFileUuid: dirFileKey,
      })
      .subscribe({
        next: (resp) => {
          this._moveEachToDir(files, dirFileKey, offset + 1);
        },
      });
  }
}
