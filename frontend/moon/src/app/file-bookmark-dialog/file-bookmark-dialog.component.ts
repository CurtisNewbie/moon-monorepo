import { Component, OnInit } from "@angular/core";
import { FileBookmark, TempFile } from "src/common/file-bookmark";
import { Toaster } from "../notification.service";
import { MatDialog } from "@angular/material/dialog";
import { DirectoryMoveFileComponent } from "../directory-move-file/directory-move-file.component";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { HttpClient } from "@angular/common/http";
import { VfolderAddFileComponent } from "../vfolder-add-file/vfolder-add-file.component";
import { guessFileIconClz, isImageByName } from "src/common/file";
import { FileType } from "src/common/file-info";
import { HostOnGalleryComponent } from "../host-on-gallery/host-on-gallery.component";

@Component({
  selector: "app-file-bookmark-dialog",
  template: `
    <h1 mat-dialog-title>Bookmarked Files</h1>
    <div mat-dialog-content>
      <ng-container *ngIf="dat.length > 0">
        <div class="d-flex flex-wrap gap-2 justify-content-end">
          <button mat-raised-button class="m-1" (click)="addToGallery()">
            Add To Gallery
          </button>
          <button mat-raised-button class="m-1" (click)="addToVirtualFolder()">
            Add To Virtual Folder
          </button>
          <button mat-raised-button class="m-1" (click)="moveToDir()">
            Move To Directory
          </button>
          <button mat-raised-button class="m-1" (click)="deleteFiles()">
            Delete
          </button>
          <button
            mat-raised-button
            class="m-1"
            (click)="clear()"
            [mat-dialog-close]="true"
          >
            Clear
          </button>
        </div>
        <div class="mt-3">
          <mat-list role="list" *ngFor="let f of dat; let i = index">
            <mat-list-item role="listitem">
              <div style="width: 100px" class="me-2">
                <img
                  style="max-height:50px; padding: 5px 0px 5px 0px;"
                  *ngIf="f.thumbnailUrl"
                  [src]="f.thumbnailUrl"
                />
                <i
                  style="max-height:50px; padding: 5px 0px 5px 0px;"
                  *ngIf="!f.thumbnailUrl"
                  [ngClass]="['bi', 'icon-button-large', guessFileIcon(f)]"
                ></i>
              </div>
              <span>{{ i + 1 }}. {{ f.name }}</span>
              <button mat-icon-button (click)="removeBookmark(f.fileKey)">
                <i class="bi bi-x icon-button-large"></i>
              </button>
            </mat-list-item>
          </mat-list>
        </div>
      </ng-container>
      <div *ngIf="!dat || dat.length < 1" class="alert alert-dark">
        <p class="mt-2 mb-3">You haven't bookmarked any file yet</p>
      </div>
    </div>
  `,
  styles: [],
})
export class FileBookmarkDialogComponent implements OnInit {
  dat: TempFile[] = [];
  guessFileIcon = guessFileIconClz;

  constructor(
    private fileBookmark: FileBookmark,
    private toaster: Toaster,
    private dialog: MatDialog,
    private http: HttpClient
  ) {}

  ngOnInit(): void {
    this.reload();
  }

  reload() {
    this.dat = [];
    for (let f of this.fileBookmark.bucket.values()) {
      this.dat.push(f);
    }
  }

  removeBookmark(fileKey: string) {
    this.fileBookmark.del(fileKey);
    this.reload();
    if (!this.dat || this.dat.length < 1) {
      this.dialog.closeAll();
    }
  }

  moveToDir() {
    this.dialog
      .open(DirectoryMoveFileComponent, {
        width: "800px",
        data: {
          files: this.dat.map((f, i) => {
            return { name: `${i + 1}. ${f.name}`, fileKey: f.fileKey };
          }),
        },
      })
      .afterClosed()
      .subscribe();
  }

  deleteFiles() {
    let msgs = [`You sure you want to delete the following files?`, ""];
    for (let s of this.dat) {
      msgs.push(s.name);
    }

    this.dialog
      .open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Delete Files",
          msg: msgs,
          isNoBtnDisplayed: true,
        },
      })
      .afterClosed()
      .subscribe((confirm) => {
        console.log(confirm);
        if (!confirm) {
          return;
        }

        let fks = [];
        for (let f of this.dat) {
          fks.push(f.fileKey);
        }
        this.http
          .post<any>(`vfm/open/api/file/delete/batch`, {
            fileKeys: fks,
          })
          .subscribe(() => {
            this.fileBookmark.clear();
            this.reload();
          });
      });
  }

  addToVirtualFolder() {
    let selected = this.dat.map((f, i) => {
      return { name: `${i + 1}. ${f.name}`, fileKey: f.fileKey };
    });

    this.dialog
      .open(VfolderAddFileComponent, {
        width: "500px",
        data: { files: selected },
      })
      .afterClosed()
      .subscribe();
  }

  addToGallery() {
    let selected = this.dat
      .filter(
        (f): boolean => isImageByName(f.name) || f.fileType == FileType.DIR
      )
      .map((f, i) => {
        return {
          name: `${i + 1}. ${f.name}`,
          fileKey: f.fileKey,
          type: f.fileType,
        };
      });

    if (!selected || selected.length < 1) {
      return;
    }

    this.dialog
      .open(HostOnGalleryComponent, {
        width: "500px",
        data: { files: selected },
      })
      .afterClosed()
      .subscribe();
  }

  clear() {
    this.fileBookmark.clear();
  }
}
