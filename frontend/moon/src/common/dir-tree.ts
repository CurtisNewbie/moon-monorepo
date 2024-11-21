import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
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
}
