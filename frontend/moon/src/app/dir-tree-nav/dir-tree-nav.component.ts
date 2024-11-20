import { NestedTreeControl } from "@angular/cdk/tree";
import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { NavType } from "../routes";
import { NavigationService } from "../navigation.service";

export interface DirTopDownTreeNode {
  fileKey?: string;
  name?: string;
  child?: DirTopDownTreeNode[];
}

@Component({
  selector: "app-dir-tree-nav",
  template: `
    <h1 mat-dialog-title>Directory Tree Navigation</h1>
    <div mat-dialog-content>
      <div>
        <mat-tree
          [dataSource]="dirTreeDataSource"
          [treeControl]="dirTreeControl"
          class="tree"
        >
          <mat-nested-tree-node *matTreeNodeDef="let node; when: treeHasChild">
            <li>
              <div class="mat-tree-node">
                <button
                  *ngIf="treeHasChild(0, node)"
                  mat-icon-button
                  matTreeNodeToggle
                >
                  <i
                    *ngIf="dirTreeControl.isExpanded(node)"
                    class="bi bi-folder2-open"
                  ></i>
                  <i
                    *ngIf="!dirTreeControl.isExpanded(node)"
                    class="bi bi-folder-fill"
                  ></i>
                </button>
                <button *ngIf="!treeHasChild(0, node)" mat-icon-button>
                  <i class="bi bi-folder2-open"></i>
                </button>
                <button mat-icon-button (click)="selectDir(node)">
                  /{{ node.name }}
                </button>
              </div>
              <ul [class.tree-invisible]="!dirTreeControl.isExpanded(node)">
                <ng-container matTreeNodeOutlet></ng-container>
              </ul>
            </li>
          </mat-nested-tree-node>
        </mat-tree>
      </div>
    </div>
  `,
  styles: [],
})
export class DirTreeNavComponent implements OnInit {
  constructor(
    private http: HttpClient,
    private snackBar: MatSnackBar,
    private nav: NavigationService
  ) {}

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  ngOnInit(): void {
    this.fetchTopDownDirTree();
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
    this.nav.navigateTo(NavType.MANAGE_FILES, [{ parentDirKey: n.fileKey }]);
  }
}
