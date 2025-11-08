import { NestedTreeControl } from "@angular/cdk/tree";
import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { DirTopDownTreeNode, DirTree } from "src/common/dir-tree";
import { Env } from "src/common/env-util";

@Component({
  selector: "app-dir-tree-nav",
  template: `
    <h3 mat-dialog-title>{{ title }}</h3>
    <div mat-dialog-content class="mt-2">
      <div>
        <mat-form-field style="width: 100%">
          <mat-label>Search Dir Name</mat-label>
          <input
            matInput
            [(ngModel)]="searchDirTreeName"
            (keyup)="onSearchDirTreeNameChanged()"
          />
        </mat-form-field>
      </div>
      <div class="mt-1 mb-1" style="max-height: 500px; overflow: scroll;">
        <mat-tree
          [dataSource]="dirTreeDataSource"
          [treeControl]="dirTreeControl"
          class="tree"
        >
          <mat-nested-tree-node
            *matTreeNodeDef="let node; when: dirTree.treeHasChild"
          >
            <li>
              <div
                class="mat-tree-node"
                *ngIf="
                  !searchDirTreeName ||
                  (searchDirTreeName && (!!node.child && node.child.length > 0) && dirTreeControl.isExpanded(node)) ||
                  (searchDirTreeName &&
                    dirTree.matchNode(node, searchDirTreeName))
                "
              >
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

                <span
                  (click)="selectDir(node)"
                  style="text-align: left; white-space: nowrap; overflow: scroll;"
                  [ngStyle]="{
                    'font-weight': dirTree.matchNode(node, searchDirTreeName)
                      ? 'bold'
                      : null,
                    'max-width': env.isMobile() ? '200px' : ''
                  }"
                  >/{{ node.name }}</span
                >
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
  searchDirTreeName: string = "";

  @Input()
  title: string = "Directory Tree Navigation";

  @Output("selected")
  selectedEmiter = new EventEmitter<DirTopDownTreeNode>();

  constructor(public env: Env, public dirTree: DirTree) {}

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  fetchTopDownDirTree() {
    this.dirTree.fetchTopDownDirTree((dat) => {
      this.dirTreeDataSource.data = [dat];
      this.dirTreeControl.dataNodes = this.dirTreeDataSource.data;
    });
  }

  ngOnInit(): void {
    this.fetchTopDownDirTree();
  }

  selectDir(node) {
    this.selectedEmiter.emit({ ...node });
  }

  onSearchDirTreeNameChanged() {
    this.dirTreeControl.collapseAll();
    if (!this.searchDirTreeName) {
      return;
    }
    this.dirTree.searchMulti(
      this.dirTreeControl,
      this.dirTreeDataSource.data,
      this.searchDirTreeName
    );
  }

  collapseAll() {
    this.dirTreeControl.collapseAll();
  }
}
