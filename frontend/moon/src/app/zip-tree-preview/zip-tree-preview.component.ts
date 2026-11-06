import { NestedTreeControl } from '@angular/cdk/tree';
import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatTreeNestedDataSource } from '@angular/material/tree';
import { guessFileIconClz, resolveSize } from 'src/common/file';
import {
  ZipTreePreview,
  ZipTreeNode,
  ZIP_TREE_TRUNCATED_BY_DEPTH,
  ZIP_TREE_TRUNCATED_BY_ENTRIES,
} from 'src/common/zip-tree';

/** What the manage-files page hands the dialog. */
export interface ZipTreePreviewData {
  /** the archive's file name, shown in the title */
  name: string;
  /** the tree the endpoint returned */
  preview: ZipTreePreview;
}

@Component({
  selector: 'app-zip-tree-preview',
  templateUrl: './zip-tree-preview.component.html',
  styleUrls: ['./zip-tree-preview.component.css'],
})
export class ZipTreePreviewComponent {
  readonly depthReason = ZIP_TREE_TRUNCATED_BY_DEPTH;
  readonly entriesReason = ZIP_TREE_TRUNCATED_BY_ENTRIES;

  treeControl = new NestedTreeControl<ZipTreeNode>((node) => node.children);
  dataSource = new MatTreeNestedDataSource<ZipTreeNode>();

  constructor(
    public dialogRef: MatDialogRef<ZipTreePreviewComponent, void>,
    @Inject(MAT_DIALOG_DATA) public data: ZipTreePreviewData
  ) {
    // the response's root is a nameless folder; the archive's own top level is
    // what the user came to see
    const roots = this.data.preview.tree?.children || [];
    this.dataSource.data = roots;
    this.treeControl.dataNodes = roots;
    // top-level folders open, anything opened from there waits for a click
    roots.forEach((node) => this.treeControl.expand(node));
  }

  hasChild(_: number, node: ZipTreeNode): boolean {
    return !!node.children && node.children.length > 0;
  }

  /** An empty folder has no children, so it reaches the leaf template. */
  iconOf(node: ZipTreeNode): string {
    return node.dir ? 'bi-folder-fill' : guessFileIconClz(node);
  }

  sizeOf(node: ZipTreeNode): string {
    return resolveSize(node.size || 0);
  }

  hasReason(reason: string): boolean {
    return !!this.data.preview.truncatedReasons?.includes(reason);
  }
}
