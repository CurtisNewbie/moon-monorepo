import { Component, Inject, OnInit, ViewChild } from "@angular/core";
import { Paging } from "src/common/paging";
import {
  canPreview,
  guessFileIconClz,
  isPdf,
  isStreamableVideo,
  isTxt,
  isWebpage,
  resolveSize,
} from "src/common/file";
import { Subscription } from "rxjs";
import { FileInfoService, TokenType } from "../file-info.service";
import { HttpClient, HttpEventType } from "@angular/common/http";
import { MediaStreamerComponent } from "../media-streamer/media-streamer.component";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { ImageViewerComponent } from "../image-viewer/image-viewer.component";
import { isEnterKey } from "src/common/condition";
import { Env } from "src/common/env-util";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ConfirmDialog } from "src/common/dialog";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface ApiDelVerFileReq {
  verFileId?: string; // Versioned File Id
}

export interface VerFileHistoryDialogData {
  verFileId?: string;
  preview: any;
}

export interface ApiListVerFileHistoryRes {
  name?: string; // file name
  fileKey?: string; // file key
  sizeInBytes?: number; // size in bytes
  uploadTime?: number; // last upload time
  thumbnail?: string; // thumbnail token
  sizeLabel?: string;
}

function preview(u, dialog, nav, fileService, isMobile, onNav = null): void {
  if (!canPreview(u.name)) {
    return;
  }

  const isStreaming = isStreamableVideo(u.name);
  fileService
    .generateFileTempToken(
      u.fileKey,
      isStreaming ? TokenType.STREAMING : TokenType.DOWNLOAD
    )
    .subscribe({
      next: (resp) => {
        const token = resp.data;

        const getDownloadUrl = () =>
          "fstore/file/raw?key=" + encodeURIComponent(token);
        const getStreamingUrl = () =>
          "fstore/file/stream?key=" + encodeURIComponent(token);

        if (isStreaming) {
          dialog.open(MediaStreamerComponent, {
            data: {
              name: u.name,
              url: getStreamingUrl(),
              token: token,
            },
          });
        } else if (isPdf(u.name)) {
          if (onNav) {
            onNav();
          }
          nav.navigateTo(NavType.PDF_VIEWER, [
            { name: u.name, url: getDownloadUrl(), uuid: u.fileKey },
          ]);
        } else if (isTxt(u.name)) {
          if (onNav) {
            onNav();
          }
          nav.navigateTo(NavType.TXT_VIEWER, [
            { name: u.name, url: getDownloadUrl(), uuid: u.fileKey },
          ]);
        } else if (isWebpage(u.name)) {
          console.log("is webpage");
          this.nav.navigateTo(NavType.WEBPAGE_VIEWER, [
            { name: u.name, url: getDownloadUrl(), uuid: u.uuid },
          ]);
        } else {
          // image
          dialog.open(ImageViewerComponent, {
            data: {
              name: u.name,
              url: getDownloadUrl(),
              isMobile: isMobile,
              rotate: false,
            },
          });
        }
      },
    });
}

@Component({
  selector: "app-ver-file-history-dialog",
  template: `
    <h1 mat-dialog-title>Versioned File History</h1>
    <div mat-dialog-content>
      <div class="mb-2">
        <p>Listing Versioned File History '{{ data.verFileId }}'</p>
      </div>

      <p class="mb-3">
        <b>Accumulated size (rough estimate): '<= {{ totalSizeLabel }}'</b>
      </p>

      <table
        mat-table
        [dataSource]="tabdata"
        class="mat-elevation-z8 mb-4"
        style="width: 100%;"
        multiTemplateDataRows
      >
        <ng-container matColumnDef="name">
          <th mat-header-cell *matHeaderCellDef><b>Name</b></th>
          <td mat-cell *matCellDef="let f">
            <span class="pl-1 pr-1">{{ f.name }}</span>
          </td>
        </ng-container>

        <ng-container matColumnDef="fileKey">
          <th mat-header-cell *matHeaderCellDef>File Key</th>
          <td mat-cell *matCellDef="let u" (click)="$event.stopPropagation()">
            {{ u.fileKey }}
          </td>
        </ng-container>

        <ng-container matColumnDef="thumbnail">
          <th mat-header-cell *matHeaderCellDef><b>Preview</b></th>
          <td mat-cell *matCellDef="let f">
            <img
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="f.thumbnail"
              [src]="f.thumbnail"
            />
            <i
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="!f.thumbnail"
              [ngClass]="['bi', 'icon-button-large', guessFileIcon(f)]"
            ></i>
          </td>
        </ng-container>

        <ng-container matColumnDef="uploadTime">
          <th mat-header-cell *matHeaderCellDef><b>Upload Time</b></th>
          <td mat-cell *matCellDef="let f">
            {{ f.uploadTime | date : "yyyy-MM-dd HH:mm:ss" }}
          </td>
        </ng-container>

        <ng-container matColumnDef="size">
          <th mat-header-cell *matHeaderCellDef><b>Size</b></th>
          <td mat-cell *matCellDef="let f">{{ f.sizeLabel }}</td>
        </ng-container>

        <ng-container matColumnDef="operate">
          <th mat-header-cell *matHeaderCellDef><b>Operate</b></th>
          <td mat-cell *matCellDef="let f">
            <button
              class="small-btn m-2"
              mat-raised-button
              (click)="$event.stopPropagation() || jumpToDownloadUrl(f.fileKey)"
            >
              <i class="bi icon-button-large bi-cloud-download"> </i>
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="columns"></tr>
        <tr
          mat-row
          *matRowDef="let row; columns: columns"
          class="element-row"
          (click)="preview(row)"
        ></tr>
      </table>

      <app-controlled-paginator
        (pageChanged)="fetch()"
      ></app-controlled-paginator>
    </div>

    <div mat-dialog-actions>
      <button mat-button [mat-dialog-close]="false" cdkFocusInitial>
        Close
      </button>
    </div>
  `,
  styles: [],
})
export class VerFileHistoryComponent implements OnInit {
  readonly columns: string[] = [
    "thumbnail",
    "name",
    "fileKey",
    "uploadTime",
    "size",
    "operate",
  ];
  tabdata: ApiListVerFileHistoryRes[] = [];
  totalSizeLabel = "unknown";

  isEnterPressed = isEnterKey;
  guessFileIcon = guessFileIconClz;
  preview = (u) => {
    preview(
      u,
      this.dialog,
      this.nav,
      this.fileService,
      this.env.isMobile(),
      () => {
        this.dialogRef.close();
      }
    );
  };

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    public dialogRef: MatDialogRef<
      VerFileHistoryComponent,
      VerFileHistoryDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: VerFileHistoryDialogData,
    private fileService: FileInfoService,
    private dialog: MatDialog,
    private nav: NavigationService,
    public env: Env,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.pagingController.setPageLimit(5);
    this.pagingController.PAGE_LIMIT_OPTIONS = [5];
    this.qryTotalSize();
  }

  qryTotalSize() {
    this.http
      .post<any>(`vfm/open/api/versioned-file/accumulated-size`, {
        verFileId: this.data.verFileId,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.totalSizeLabel = resolveSize(resp.data.sizeInBytes);
        },
      });
  }

  fetch() {
    this.http
      .post<any>(`vfm/open/api/versioned-file/history`, {
        paging: this.pagingController.paging,
        verFileId: this.data.verFileId,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          if (!resp.data.payload) {
            resp.data.payload = [];
          }
          this.tabdata = resp.data.payload;
          this.pagingController.onTotalChanged(resp.data.paging);
          for (let f of this.tabdata) {
            if (f.thumbnail) {
              f.thumbnail =
                "fstore/file/raw?key=" + encodeURIComponent(f.thumbnail);
            }
            f.sizeLabel = resolveSize(f.sizeInBytes);
          }
        },
      });
  }

  jumpToDownloadUrl(fileKey: string): void {
    this.fileService.jumpToDownloadUrl(fileKey);
  }
}

export interface ApiListVerFileReq {
  paging?: Paging;
  name?: string; // file name
}

export interface ApiListVerFileRes {
  verFileId?: string; // versioned file id
  name?: string; // file name
  fileKey?: string; // file key
  sizeInBytes?: number; // size in bytes
  uploadTime?: number; // last upload time
  createTime?: number; // create time of the versioned file record
  updateTime?: number; // Update time of the versioned file record
  thumbnail?: string; // thumbnail token

  // ------
  sizeLabel?: string;
}

@Component({
  selector: "app-versioned-file",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Versioned Files</h3>
    </div>

    <div
      *ngIf="expandUploadPanel"
      class="container-fluid bootstrap p-3 shadow rounded mt-3 border"
    >
      <ng-container *ngIf="!updateVerFileId">
        <h4 class="mt-3 mb-3">Create Versioned File</h4>
      </ng-container>
      <ng-container *ngIf="updateVerFileId">
        <h4 class="mt-3 mb-3">
          Update Versioned File {{ updateVerFileId }} {{ updateVerFileName }}
        </h4>
      </ng-container>

      <mat-form-field style="width: 100%;" class="mt-1 mb-1">
        <mat-label>Name</mat-label>
        <input
          matInput
          type="text"
          [(ngModel)]="uploadFileName"
          [disabled]="isUploading"
        />
        <button
          *ngIf="uploadFileName"
          matSuffix
          aria-label="Clear"
          (click)="uploadFileName = ''"
          class="btn-close"
          [disabled]="isUploading"
        ></button>
      </mat-form-field>

      <div class="input-group input-group-lg mt-1 mb-1">
        <input
          type="file"
          class="form-control darkmode"
          #uploadFileInput
          (change)="onFileSelected($event.target.files)"
          aria-describedby="basic-addon1"
          [disabled]="isUploading"
        />
      </div>

      <div class="mt-3 mb-2">
        <div class="row row-cols-lg-auto g-3 align-items-center">
          <div class="col">
            <button
              class="ml-2 mr-2"
              mat-raised-button
              (click)="upload()"
              [disabled]="isUploading"
              *ngIf="!isUploading"
            >
              Upload
            </button>
            <button
              class="ml-2 mr-2"
              *ngIf="isUploading"
              mat-raised-button
              (click)="cancelFileUpload()"
            >
              Cancel
            </button>
          </div>
          <div class="col">
            <small style="color: cadetblue;" *ngIf="progress != null"
              >Progress: {{ progress }}</small
            >
          </div>
        </div>
      </div>
      <!-- upload param end -->
    </div>

    <mat-form-field style="width: 100%;" class="mb-1 mt-3">
      <mat-label>Versioned File Name:</mat-label>
      <input matInput type="text" [(ngModel)]="searchName" />
    </mat-form-field>

    <div class="d-grid gap-2 d-flex justify-content-end mb-3">
      <button
        mat-icon-button
        class="m-1 icon-button-large"
        [class.status-green]="expandUploadPanel"
        (click)="toggleUploadPanel()"
      >
        <i class="bi bi-cloud-upload"></i>
      </button>
      <button mat-icon-button class="m-1 icon-button-large" (click)="fetch()">
        <i class="bi bi-arrow-clockwise"></i>
      </button>
      <button mat-icon-button class="m-1 icon-button-large" (click)="reset()">
        <i class="bi bi-slash-circle"></i>
      </button>
    </div>

    <div class="mt-3 mb-2" style="overflow: auto;">
      <table mat-table [dataSource]="tabdat" style="width: 100%;">
        <ng-container matColumnDef="name">
          <th mat-header-cell *matHeaderCellDef><b>Name</b></th>
          <td mat-cell *matCellDef="let f">
            <span class="pl-1 pr-1">{{ f.name }} </span>
          </td>
        </ng-container>

        <ng-container matColumnDef="verFileId">
          <th mat-header-cell *matHeaderCellDef>Versioned File Id</th>
          <td mat-cell *matCellDef="let u">
            {{ u.verFileId }}
          </td>
        </ng-container>

        <ng-container matColumnDef="thumbnail">
          <th mat-header-cell *matHeaderCellDef><b>Preview</b></th>
          <td mat-cell *matCellDef="let f">
            <img
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="f.thumbnail"
              [src]="f.thumbnail"
            />
            <i
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="!f.thumbnail"
              [ngClass]="['bi', 'icon-button-large', guessFileIcon(f)]"
            ></i>
          </td>
        </ng-container>

        <ng-container matColumnDef="uploadTime">
          <th mat-header-cell *matHeaderCellDef><b>Upload Time</b></th>
          <td mat-cell *matCellDef="let f">
            {{ f.uploadTime | date : "yyyy-MM-dd HH:mm:ss" }}
          </td>
        </ng-container>

        <ng-container matColumnDef="size">
          <th mat-header-cell *matHeaderCellDef><b>Size</b></th>
          <td mat-cell *matCellDef="let f">{{ f.sizeLabel }}</td>
        </ng-container>

        <ng-container matColumnDef="operate">
          <th mat-header-cell *matHeaderCellDef><b>Operate</b></th>
          <td mat-cell *matCellDef="let f">
            <button
              mat-raised-button
              class="m-2"
              (click)="$event.stopPropagation() || selectVerFile(f)"
            >
              Update
            </button>
            <button
              mat-raised-button
              class="m-2"
              (click)="$event.stopPropagation() || showVerFileHistory(f)"
            >
              History
            </button>
            <button
              mat-raised-button
              class="m-2"
              (click)="$event.stopPropagation() || deleteVerFile(f)"
            >
              Delete
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="tabcol"></tr>
        <tr
          mat-row
          *matRowDef="let row; columns: tabcol"
          (click)="preview(row)"
        ></tr>
      </table>
    </div>

    <app-controlled-paginator
      (pageChanged)="fetch()"
    ></app-controlled-paginator>
  `,
  styles: [],
})
export class VersionedFileComponent implements OnInit {
  expandUploadPanel = false;
  tabdat: ApiListVerFileRes[] = [];
  tabcol = ["thumbnail", "verFileId", "name", "uploadTime", "size", "operate"];
  searchName = "";
  updateVerFileId = "";
  updateVerFileName = "";

  progress: string = null;
  isUploading: boolean = false;
  uploadFileName: string = null;
  uploadFile: File = null;
  uploadSub: Subscription = null;

  guessFileIcon = guessFileIconClz;
  preview = (u) => {
    preview(u, this.dialog, this.nav, this.fileService, this.env.isMobile());
  };

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private http: HttpClient,
    private fileService: FileInfoService,
    private dialog: MatDialog,
    private nav: NavigationService,
    public env: Env,
    private snackBar: MatSnackBar,
    private confirm: ConfirmDialog
  ) {}

  ngOnInit(): void {}

  fetch() {
    let req: ApiListVerFileReq = {
      paging: this.pagingController.paging,
      name: this.searchName,
    };
    this.http.post<any>(`vfm/open/api/versioned-file/list`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        if (!resp.data.payload) {
          resp.data.payload = [];
        }
        this.tabdat = resp.data.payload;
        this.pagingController.onTotalChanged(resp.data.paging);
        for (let f of this.tabdat) {
          if (f.thumbnail) {
            f.thumbnail =
              "fstore/file/raw?key=" + encodeURIComponent(f.thumbnail);
          }
          f.sizeLabel = resolveSize(f.sizeInBytes);
        }
      },
    });
  }

  reset() {
    this.pagingController.firstPage();
    this.searchName = "";
  }

  onFileSelected(files: File[]): void {
    if (this.isUploading) return;

    if (files.length < 1) {
      return;
    }

    this.uploadFile = files[0];
    this.uploadFileName = this.uploadFile.name;
  }

  cancelFileUpload(): void {
    if (!this.isUploading) return;

    if (this.uploadSub != null && !this.uploadSub.closed) {
      this.uploadSub.unsubscribe();
      return;
    }

    this.progress = null;
    this.isUploading = false;
    this.uploadFile = null;
    this.uploadFileName = "";
    this.snackBar.open("File uploading cancelled", "ok", { duration: 3000 });
  }

  upload() {
    if (!this.uploadFile) {
      this.snackBar.open("Please select file to upload", "ok", {
        duration: 3000,
      });
      return;
    }

    const abortUpload = () => {
      this.progress = null;
      this.isUploading = false;
      this.uploadFile = null;
      this.uploadFileName = "";
      this.snackBar.open("Failed to upload file", "ok", { duration: 3000 });
    };

    this.uploadSub = this.fileService
      .uploadToMiniFstore({
        fileName: this.uploadFileName,
        files: [this.uploadFile],
      })
      .subscribe({
        next: (event) => {
          if (event.type === HttpEventType.UploadProgress) {
            this.updateProgress(this.uploadFileName, event.loaded, event.total);
          }

          if (event.type == HttpEventType.Response) {
            let fstoreRes = event.body;
            if (fstoreRes.error) {
              abortUpload();
              return;
            }
            let sub = null;
            if (this.updateVerFileId) {
              sub = this.http.post(`vfm/open/api/versioned-file/update`, {
                filename: this.uploadFileName,
                fstoreFileId: fstoreRes.data,
                verFileId: this.updateVerFileId,
              });
            } else {
              sub = this.http.post(`vfm/open/api/versioned-file/create`, {
                filename: this.uploadFileName,
                fstoreFileId: fstoreRes.data,
              });
            }

            sub.subscribe({
              next: (resp) => {
                if (resp.error) {
                  this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                  return;
                }
              },
              complete: () => {
                this.progress = null;
                this.isUploading = false;
                this.uploadFile = null;
                this.uploadFileName = "";
                this.expandUploadPanel = false;
                this.updateVerFileId = "";
                this.updateVerFileName = "";
                this.fetch();
              },
              error: () => {
                abortUpload();
              },
            });
          }
        },
        error: () => {
          abortUpload();
        },
      });
  }

  updateProgress(filename: string, loaded: number, total: number) {
    let p = Math.round((100 * loaded) / total).toFixed(2);
    let ps;
    if (p == "100.00") ps = `Processing '${filename}' ... `;
    else ps = `Uploading ${filename} ${p}% `;
    this.progress = ps;
  }

  toggleUploadPanel() {
    if (this.isUploading) {
      return;
    }
    this.expandUploadPanel = !this.expandUploadPanel;
    this.updateVerFileId = "";
    this.updateVerFileName = "";
  }

  selectVerFile(f: ApiListVerFileRes) {
    if (this.expandUploadPanel && this.updateVerFileId == f.verFileId) {
      this.toggleUploadPanel();
      return;
    }

    this.expandUploadPanel = true;
    this.updateVerFileId = f.verFileId;
    this.updateVerFileName = f.name;
  }

  showVerFileHistory(f: ApiListVerFileRes) {
    const dialogRef: MatDialogRef<VerFileHistoryComponent, boolean> =
      this.dialog.open(VerFileHistoryComponent, {
        width: "80vw",
        data: { verFileId: f.verFileId, preview: () => this.preview },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      // do nothing
    });
  }

  deleteVerFile(f: ApiListVerFileRes) {
    this.confirm.show(
      `Delete ${f.name}?`,
      [
        `Are you sure you want to delete ${f.name}?`,
        "All snapshots are deleted as well.",
      ],
      () => {
        let req: ApiDelVerFileReq = { verFileId: f.verFileId };
        this.http
          .post<any>(`/vfm/open/api/versioned-file/delete`, req)
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.fetch();
            },
            error: (err) => {
              console.log(err);
              this.snackBar.open("Request failed, unknown error", "ok", {
                duration: 3000,
              });
            },
          });
      }
    );
  }
}
