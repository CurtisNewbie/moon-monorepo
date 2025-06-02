import { NestedTreeControl } from "@angular/cdk/tree";
import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { NavigationService } from "src/app/navigation.service";
import { NavType } from "src/app/routes";

export interface DirTopDownTreeNode {
  fileKey?: string;
  name?: string;
  child?: DirTopDownTreeNode[];
}

@Injectable({
  providedIn: "root",
})
export class DirTree {
  constructor(
    private http: HttpClient,
    private snackBar: MatSnackBar,
    private nav: NavigationService
  ) {}

  newDirTreeControl(): NestedTreeControl<DirTopDownTreeNode> {
    return new NestedTreeControl<DirTopDownTreeNode>((node) => node.child);
  }

  newDirTreeDataSource(): MatTreeNestedDataSource<DirTopDownTreeNode> {
    return new MatTreeNestedDataSource<DirTopDownTreeNode>();
  }

  fetchTopDownDirTree(onDirTreeFetched) {
    this.http.get<any>(`/vfm/open/api/file/dir/top-down-tree`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: DirTopDownTreeNode = resp.data;
        onDirTreeFetched(dat);
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

  goToFile(n) {
    this.nav.navigateTo(NavType.MANAGE_FILES, [{ parentDirKey: n.fileKey }]);
  }

  search(
    dtc: NestedTreeControl<DirTopDownTreeNode>,
    root: DirTopDownTreeNode,
    kw: string
  ): boolean {
    if (!root) {
      return false;
    }
    if (this.matchNode(root, kw)) {
      console.log("found, ", root);
      return true;
    }
    let rootFound = false;
    if (root.child && root.child.length > 0) {
      for (let c of root.child) {
        if (this.search(dtc, c, kw)) {
          console.log("expand, ", root);
          dtc.expand(root);
          rootFound = true;
        }
      }
    }
    return rootFound;
  }

  searchMulti(
    dtc: NestedTreeControl<DirTopDownTreeNode>,
    roots: DirTopDownTreeNode[],
    kw: string
  ) {
    console.log("roots, ", roots);
    for (let r of roots) {
      if (this.search(dtc, r, kw)) {
        console.log("root found, ", r);
      }
    }
  }

  matchNode(n: DirTopDownTreeNode, kw: string): boolean {
    return n.name && n.name.toLowerCase().indexOf(kw.toLowerCase()) > -1;
  }
}
