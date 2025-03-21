import { Component, Inject, OnInit } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NestedTreeControl } from "@angular/cdk/tree";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { DirTopDownTreeNode, DirTree } from "src/common/dir-tree";

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
  /** name of dir that we may move file into */
  moveIntoDirName: string = null;
  moveIntoDirKey: string = null;

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  constructor(
    public dialogRef: MatDialogRef<DirectoryMoveFileComponent, Data>,
    @Inject(MAT_DIALOG_DATA) public dat: Data,
    private http: HttpClient,
    private snackBar: MatSnackBar,
    public dirTree: DirTree
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
      .post<any>(`/vfm/open/api/file/batch-move-to-dir`, { instructions: reqs })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
        },
      });
  }

  fetchTopDownDirTree() {
    this.dirTree.fetchTopDownDirTree((dat) => {
      this.dirTreeDataSource.data = [dat];
      this.dirTreeControl.dataNodes = this.dirTreeDataSource.data;
    });
  }

  selectDir(n) {
    this.moveIntoDirKey = n.fileKey;
    this.moveIntoDirName = n.name;
    this.snackBar.open(`Selected directory /${n.name}`, "ok", {
      duration: 1500,
    });
    this.dirTreeControl.collapseAll();
  }
}
