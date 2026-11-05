import { NestedTreeControl } from "@angular/cdk/tree";
import {
  AfterViewChecked,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild,
} from "@angular/core";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { DirTopDownTreeNode, DirTree } from "src/common/dir-tree";

/** Where an applied selection came from. */
interface ApplySelectionOptions {
  /**
   * An input applies once per distinct key and defers the emit to a microtask.
   * A click always applies, even on the node that is already selected, and
   * emits synchronously.
   */
  source?: "input" | "click";
}

@Component({
  selector: "app-dir-tree-nav",
  template: `
    <h3 mat-dialog-title>{{ 'dir-tree-nav' | trl:title }}</h3>
    <div mat-dialog-content class="mt-2">
      <div>
        <mat-form-field style="width: 100%">
          <mat-label>{{ 'dir-tree-nav' | trl:'searchDirName' }}</mat-label>
          <input
            matInput
            [(ngModel)]="searchDirTreeName"
            (keyup)="onSearchDirTreeNameChanged()"
          />
        </mat-form-field>
      </div>
      <div
        #treeScrollBox
        class="mt-1 mb-1 dir-tree-box"
        style="max-height: 500px; overflow: scroll;"
      >
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
                  (searchDirTreeName &&
                    !!node.child &&
                    node.child.length > 0 &&
                    dirTreeControl.isExpanded(node)) ||
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
                  class="dir-tree-name"
                  [class.dir-tree-name-selected]="node === selectedNode"
                  [attr.aria-selected]="node === selectedNode ? 'true' : null"
                  style="text-align: left; white-space: nowrap; overflow: scroll;"
                  [ngStyle]="{
                    'font-weight': dirTree.matchNode(node, searchDirTreeName)
                      ? 'bold'
                      : null
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
  styleUrls: ["./dir-tree-nav.component.css"],
})
export class DirTreeNavComponent implements OnInit, AfterViewChecked {
  searchDirTreeName: string = "";

  @Input()
  title: string = "defaultTitle";

  // FileKey to select. Applied immediately once the tree is loaded, otherwise
  // kept pending until the fetch reports back.
  private selectedFileKeyValue?: string;
  @Input()
  set selectedFileKey(v: string | undefined) {
    this.selectedFileKeyValue = v || undefined;
    if (this.dirTreeDataSource.data && this.dirTreeDataSource.data.length > 0) {
      this.applySelection(this.selectedFileKeyValue);
    }
  }

  selectedNode?: DirTopDownTreeNode;
  private appliedFileKey?: string;
  private revealSelected = false;

  @ViewChild("treeScrollBox", { static: true })
  private treeScrollBox?: ElementRef<HTMLElement>;

  // Bumped whenever a selection is applied, so an emit deferred by an earlier
  // input application is dropped once a later input or click supersedes it.
  private emitToken = 0;

  @Output("selected")
  selectedEmiter = new EventEmitter<DirTopDownTreeNode>();

  @Output("treeLoaded")
  treeLoadedEmiter = new EventEmitter<void>();

  constructor(public dirTree: DirTree) {}

  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();

  fetchTopDownDirTree() {
    this.dirTree.fetchTopDownDirTree((dat) => {
      this.dirTreeDataSource.data = [dat];
      this.dirTreeControl.dataNodes = this.dirTreeDataSource.data;
      if (this.selectedFileKeyValue) {
        this.applySelection(this.selectedFileKeyValue);
      }
      this.treeLoadedEmiter.emit();
    });
  }

  ngOnInit(): void {
    this.fetchTopDownDirTree();
  }

  ngAfterViewChecked(): void {
    if (!this.revealSelected) {
      return;
    }
    this.revealSelected = false;
    const box = this.treeScrollBox?.nativeElement;
    const el = box?.querySelector<HTMLElement>(".dir-tree-name-selected");
    if (!box || !el) {
      return;
    }
    // Scroll only the tree's own box, never an ancestor or the page. Centre the
    // label by moving the box the distance between the two centres.
    const boxRect = box.getBoundingClientRect();
    const elRect = el.getBoundingClientRect();
    box.scrollTop +=
      elRect.top + elRect.height / 2 - (boxRect.top + boxRect.height / 2);
  }

  /**
   * Locates the node, expands its ancestors, marks it, and emits it. Input
   * applications keep the once-per-distinct-key rule and defer the emit to a
   * microtask: the input is bound inside a change-detection pass, and emitting
   * there would mutate the parent's bound fields mid-pass. A click always
   * applies, even to the already selected node, and emits synchronously.
   */
  private applySelection(
    fileKey: string | undefined,
    options: ApplySelectionOptions = {}
  ): void {
    const fromClick = options.source === "click";
    if (!fromClick && fileKey && fileKey === this.appliedFileKey) {
      return;
    }
    this.emitToken++;
    this.appliedFileKey = fileKey;
    this.selectedNode = undefined;
    if (!fileKey) {
      return;
    }
    const path = this.findPath(fileKey);
    if (!path) {
      return;
    }
    this.selectedNode = path[path.length - 1];
    if (!fromClick) {
      this.prefillSearch(this.selectedNode.name);
    }
    this.expandPathAndReveal(path);
    const node = { ...this.selectedNode };
    if (fromClick) {
      this.selectedEmiter.emit(node);
      return;
    }
    const token = this.emitToken;
    Promise.resolve().then(() => {
      if (token === this.emitToken) {
        this.selectedEmiter.emit(node);
      }
    });
  }

  /**
   * Fills the search box with the located directory's name and filters the tree
   * to its matches, so an input-applied selection opens on its own branch. The
   * search expands the path to every match, including the selected one.
   */
  private prefillSearch(name: string | undefined): void {
    if (!name) {
      return;
    }
    this.searchDirTreeName = name;
    this.onSearchDirTreeNameChanged();
  }

  /** Expands every ancestor on the path, then schedules the reveal. */
  private expandPathAndReveal(path: DirTopDownTreeNode[]): void {
    if (!this.searchDirTreeName) {
      for (const ancestor of path.slice(0, -1)) {
        this.dirTreeControl.expand(ancestor);
      }
    }
    this.revealSelected = true;
  }

  /** The root-to-node chain, the node itself last, or undefined when absent. */
  private findPath(fileKey: string): DirTopDownTreeNode[] | undefined {
    const visit = (
      node: DirTopDownTreeNode,
      path: DirTopDownTreeNode[]
    ): DirTopDownTreeNode[] | undefined => {
      const current = path.concat(node);
      if (node.fileKey === fileKey) {
        return current;
      }
      for (const child of node.child || []) {
        const found = visit(child, current);
        if (found) {
          return found;
        }
      }
      return undefined;
    };
    for (const root of this.dirTreeDataSource.data) {
      const found = visit(root, []);
      if (found) {
        return found;
      }
    }
    return undefined;
  }

  selectDir(node: DirTopDownTreeNode) {
    if (node.fileKey) {
      this.applySelection(node.fileKey, { source: "click" });
    }
  }

  onSearchDirTreeNameChanged() {
    this.dirTreeControl.collapseAll();
    if (!this.searchDirTreeName) {
      this.restoreSelectedPath();
      return;
    }
    this.dirTree.searchMulti(
      this.dirTreeControl,
      this.dirTreeDataSource.data,
      this.searchDirTreeName
    );
  }

  /**
   * Re-expands the path to the selected node after a search is cleared, so the
   * highlight is visible again. The selection itself is unchanged, so nothing
   * is emitted and the tree is not fetched again.
   */
  private restoreSelectedPath(): void {
    if (!this.selectedNode || !this.selectedNode.fileKey) {
      return;
    }
    const path = this.findPath(this.selectedNode.fileKey);
    if (!path) {
      return;
    }
    this.expandPathAndReveal(path);
  }
}
