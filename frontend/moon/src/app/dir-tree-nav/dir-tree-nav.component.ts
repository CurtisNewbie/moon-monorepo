import { NestedTreeControl } from "@angular/cdk/tree";
import { Component, OnInit } from "@angular/core";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { DirTopDownTreeNode, DirTree } from "src/common/dir-tree";

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
          <mat-nested-tree-node
            *matTreeNodeDef="let node; when: dirTree.treeHasChild"
          >
            <li>
              <div class="mat-tree-node">
                <button
                  *ngIf="dirTree.treeHasChild(0, node)"
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
                <button *ngIf="!dirTree.treeHasChild(0, node)" mat-icon-button>
                  <i class="bi bi-folder2-open"></i>
                </button>
                <button mat-icon-button (click)="dirTree.goToFile(node)">
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
  constructor(public dirTree: DirTree) {}

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  fetchTopDownDirTree() {
    this.dirTree.fetchTopDownDirTree((dat) => {
      this.dirTreeDataSource.data = [dat];
      this.dirTreeControl.dataNodes = this.dirTreeDataSource.data;
      // this.dirTreeControl.expandAll();
    });
  }

  ngOnInit(): void {
    this.fetchTopDownDirTree();
  }
}
