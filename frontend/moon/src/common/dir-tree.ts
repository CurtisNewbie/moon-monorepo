import { HighContrastModeDetector } from "@angular/cdk/a11y";
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
  hidden?: boolean;
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

    let rootFound = false;
    if (this.matchNode(root, kw)) {
      console.log("found, ", root);
      rootFound = true;
    }

    if (!rootFound) {
      root.hidden = true;
    }

    if (root.child && root.child.length > 0) {
      for (let c of root.child) {
        c.hidden = true;
        if (this.search(dtc, c, kw)) {
          console.log("expand, ", root);
          dtc.expand(root);
          rootFound = true;
          c.hidden = false;
        }
      }
    }

    if (root.hidden && rootFound) {
      root.hidden = false;
    }

    return rootFound;
  }

  resetHidden(r: DirTopDownTreeNode) {
    if (!r) {
      return;
    }
    r.hidden = false;
    if (r.child && r.child.length > 0) {
      for (let c of r.child) {
        this.resetHidden(c);
      }
    }
  }

  resetHiddenNodes(roots: DirTopDownTreeNode[]) {
    for (let r of roots) {
      this.resetHidden(r);
    }
  }

  searchMulti(
    dtc: NestedTreeControl<DirTopDownTreeNode>,
    roots: DirTopDownTreeNode[],
    kw: string
  ) {
    this.resetHiddenNodes(roots);

    console.log("roots, ", roots);
    for (let r of roots) {
      if (this.search(dtc, r, kw)) {
        console.log("root found, ", r);
      }
    }
  }

  matchNode(n: DirTopDownTreeNode, kw: string): boolean {
    return kw && n.name && n.name.toLowerCase().indexOf(kw.toLowerCase()) > -1;
  }
}
