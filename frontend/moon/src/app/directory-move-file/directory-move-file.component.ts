import { Component, Inject, OnInit } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { DirBrief } from "src/common/file-info";
import { filterAlike } from "src/common/select-util";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NestedTreeControl } from "@angular/cdk/tree";
import { MatTreeNestedDataSource } from "@angular/material/tree";

export interface DirTopDownTreeNode {
  fileKey?: string;
  name?: string;
  child?: DirTopDownTreeNode[];
}

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
  moveIntoDirKey: string = null;

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  onMoveIntoDirNameChanged = () =>
    (this.autoCompMoveIntoDirs = filterAlike(
      this.dirBriefList.map((v) => v.name),
      this.moveIntoDirName
    ));

  constructor(
    public dialogRef: MatDialogRef<DirectoryMoveFileComponent, Data>,
    @Inject(MAT_DIALOG_DATA) public dat: Data,
    private http: HttpClient,
    private snackBar: MatSnackBar
  ) {
    this.dirTreeDataSource.data = [];
  }

  ngOnInit(): void {
    this.fetchTopDownDirTree();
  }

  moveToDir() {
    const dk = this.moveIntoDirKey;
    if (dk == null) {
      this.snackBar.open("Please select directory", "ok", { duration: 3000 });
      return;
    }
    this.batchMoveEachToDir(this.dat.files, dk);
  }

  private batchMoveEachToDir(files, dirFileKey: string) {
    let reqs = [];
    for (let f of files) {
      reqs.push({ uuid: f.fileKey, parentFileUuid: dirFileKey });
    }
    this.http
      .post(`/vfm/open/api/file/batch-move-to-dir`, { instructions: reqs })
      .subscribe();
  }

  fetchTopDownDirTree() {
    this.http.get<any>(`/vfm/open/api/file/dir/top-down-tree`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: DirTopDownTreeNode = resp.data;
        this.dirTreeDataSource.data = [dat];
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  treeHasChild(_: number, node: DirTopDownTreeNode) {
    return !!node.child && node.child.length > 0;
  }

  selectDir(n) {
    this.moveIntoDirKey = n.fileKey;
    this.moveIntoDirName = n.name;
    this.snackBar.open(`Selected directory /${n.name}`, "ok", {
      duration: 1500,
    });
  }
}
