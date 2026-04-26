import { Component, OnInit } from "@angular/core";
import { FileBookmark, TempFile } from "src/common/file-bookmark";
import { MatDialog } from "@angular/material/dialog";
import { DirectoryMoveFileComponent } from "../directory-move-file/directory-move-file.component";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { HttpClient } from "@angular/common/http";
import { VfolderAddFileComponent } from "../vfolder-add-file/vfolder-add-file.component";
import { guessFileIconClz, isImageByName } from "src/common/file";
import { FileType } from "src/common/file-info";
import { HostOnGalleryComponent } from "../host-on-gallery/host-on-gallery.component";
import { MatSnackBar } from "@angular/material/snack-bar";
import { I18n } from "../i18n.service";

@Component({
  selector: "app-file-tray",
  template: `
    <div class="d-flex align-items-center justify-content-between mb-4">
      <h2 mat-dialog-title class="m-0">{{'file-tray' | trl:'fileTray'}}</h2>
      <span class="badge bg-secondary" *ngIf="dat.length > 0">{{'file-tray' | trl:'itemCount':'count':dat.length}}</span>
    </div>

    <mat-dialog-content class="p-0 pt-2">
      <div *ngIf="dat.length > 0" class="d-flex flex-wrap justify-content-between align-items-center gap-2 mb-3 pb-2 border-bottom">
        <div class="d-flex flex-wrap gap-2 ms-2">
          <button mat-flat-button color="warn" (click)="deleteFiles()">{{'file-tray' | trl:'delete'}}</button>
          <button mat-stroked-button (click)="clear()">{{'file-tray' | trl:'clear'}}</button>
        </div>
        <div class="d-flex flex-wrap gap-2 me-2">
          <button mat-stroked-button (click)="addToGallery()">{{'file-tray' | trl:'gallery'}}</button>
          <button mat-stroked-button (click)="addToVirtualFolder()">{{'file-tray' | trl:'addToFolder'}}</button>
          <button mat-stroked-button (click)="moveToDir()">{{'file-tray' | trl:'move'}}</button>
        </div>
      </div>

      <div *ngIf="dat.length > 0" class="basket-list">
        <div *ngFor="let f of dat; let i = index" class="basket-item d-flex align-items-center gap-3 p-2 rounded">
          <div class="basket-thumb d-flex align-items-center justify-content-center">
            <img *ngIf="f.thumbnailUrl" [src]="f.thumbnailUrl" class="basket-thumb-img" [appImageTooltip]="f.name" [imageUrl]="f.thumbnailUrl"/>
            <i *ngIf="!f.thumbnailUrl" [ngClass]="['bi', guessFileIcon(f)]" class="basket-thumb-icon" [matTooltip]="f.name"></i>
          </div>
          <div class="flex-grow-1 min-w-0">
            <div class="basket-name text-truncate" [title]="f.name">{{ f.name }}</div>
            <div class="basket-meta text-muted">#{{ i + 1 }}</div>
          </div>
          <button mat-icon-button class="basket-remove" (click)="removeBookmark(f.fileKey)">
            <i class="bi bi-x"></i>
          </button>
        </div>
      </div>

      <div *ngIf="!dat || dat.length < 1" class="basket-empty d-flex flex-column align-items-center justify-content-center py-5">
        <i class="bi bi-basket3 basket-empty-icon mb-3"></i>
        <p class="basket-empty-text text-muted mb-0">{{'file-tray' | trl:'trayIsEmpty'}}</p>
        <p class="basket-empty-hint text-muted small mt-1">{{'file-tray' | trl:'bookmarkFilesToCollect'}}</p>
      </div>
    </mat-dialog-content>

    <mat-dialog-actions align="end" class="mt-3 mb-2 me-2">
      <button mat-stroked-button [mat-dialog-close]="true">{{'file-tray' | trl:'close'}}</button>
    </mat-dialog-actions>
  `,
  styles: [`
    :host { display: block; font-size: 12px; }
    h2[mat-dialog-title] { font-size: 1.4rem; font-weight: 500; }
    .badge { font-size: 0.85rem; padding: 0.4em 0.6em; }
    .basket-list { max-height: 320px; overflow-y: auto; }
    .basket-item { transition: background-color 0.15s ease; border-bottom: 1px solid rgba(0,0,0,0.06); }
    .basket-item:last-child { border-bottom: none; }
    .basket-item:hover { background-color: rgba(0,0,0,0.03); }
    .basket-thumb { width: 48px; height: 48px; flex-shrink: 0; border-radius: 4px; background-color: rgba(0,0,0,0.04); overflow: hidden; }
    .basket-thumb-img { max-width: 100%; max-height: 100%; object-fit: cover; }
    .basket-thumb-icon { font-size: 1.5rem; color: #6c757d; }
    .basket-name { font-size: 1rem; font-weight: 400; color: #fff; line-height: 1.4; }
    .basket-meta { font-size: 0.85rem; margin-top: 2px; color: #e0e0e0; }
    .basket-remove { opacity: 0.6; transition: opacity 0.15s ease; width: 32px; height: 32px; line-height: 32px; }
    .basket-item:hover .basket-remove { opacity: 1; }
    .basket-remove i { font-size: 1rem; }
    .basket-empty { min-height: 180px; }
    .basket-empty-icon { font-size: 3.5rem; color: #adb5bd; }
    .basket-empty-text { font-size: 1.1rem; }
    .basket-empty-hint { font-size: 0.9rem; }
    mat-dialog-content { max-height: 60vh; }
    mat-dialog-actions { padding: 0; min-height: unset; }
    button[mat-flat-button], button[mat-stroked-button] { font-size: 0.9rem; line-height: 2rem; padding: 0 1rem; }
    button[mat-flat-button] i, button[mat-stroked-button] i { font-size: 0.9rem; }
    .border-bottom { border-bottom: 1px solid rgba(0,0,0,0.12) !important; }
    .gap-1 { gap: 0.25rem !important; }
    .gap-2 { gap: 0.5rem !important; }
    .gap-3 { gap: 1rem !important; }
    .min-w-0 { min-width: 0 !important; }
  `],
})
export class FileTrayComponent implements OnInit {
  dat: TempFile[] = [];
  guessFileIcon = guessFileIconClz;

  constructor(
    private fileBookmark: FileBookmark,
    private dialog: MatDialog,
    private http: HttpClient,
    private snackBar: MatSnackBar,
    private i18n: I18n
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
      .subscribe((moved) => {
        if (!moved) {
          return;
        }
        this.fileBookmark.clear();
        this.reload();
      });
  }

  deleteFiles() {
    let msgs = [this.i18n.trl("file-tray", "confirmDeleteFiles", "count", this.dat.length), ""];
    // for (let s of this.dat) {
    //   msgs.push(s.name);
    // }

    this.dialog
      .open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: this.i18n.trl("file-tray", "deleteFilesTitle"),
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
          .subscribe((resp) => {
            if (resp.error) {
              this.snackBar.open(resp.msg, "ok", { duration: 6000 });
              return;
            }
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
      .subscribe((added) => {
        if (added) {
          this.fileBookmark.clear();
          this.reload();
        }
      });
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
      .subscribe((added) => {
        if (added) {
          this.fileBookmark.clear();
          this.reload();
        }
      });
  }

  clear() {
    this.fileBookmark.clear();
    this.reload();
  }
}
