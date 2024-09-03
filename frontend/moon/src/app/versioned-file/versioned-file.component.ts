import { Component, Inject, OnInit } from "@angular/core";
import { Paging, PagingController } from "src/common/paging";
import { environment } from "src/environments/environment";
import {
  canPreview,
  guessFileThumbnail,
  isPdf,
  isStreamableVideo,
  isTxt,
  isWebpage,
  resolveSize,
} from "src/common/file";
import { Subscription } from "rxjs";
import { FileInfoService, TokenType } from "../file-info.service";
import { HttpClient, HttpEventType } from "@angular/common/http";
import { Toaster } from "../notification.service";
import { MediaStreamerComponent } from "../media-streamer/media-streamer.component";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { ImageViewerComponent } from "../image-viewer/image-viewer.component";
import { isMobile } from "src/common/env-util";
import { isEnterKey } from "src/common/condition";

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
          environment.fstore + "/file/raw?key=" + encodeURIComponent(token);
        const getStreamingUrl = () =>
          environment.fstore + "/file/stream?key=" + encodeURIComponent(token);

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
            console.log("is webpage")
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
          <td mat-cell *matCellDef="let f" (click)="preview(f)">
            <img
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="f.thumbnail"
              [src]="f.thumbnail"
            />
            <img
              style="max-height:40px; padding: 5px 0px 5px 0px;"
              *ngIf="!f.thumbnail"
              [src]="guessFileThumbnail(f)"
            />
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
              <img style="max-height:20px;" src="../assets/download.png" />
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="columns"></tr>
        <tr
          mat-row
          *matRowDef="let row; columns: columns"
          class="element-row"
        ></tr>
      </table>

      <app-controlled-paginator
        (controllerReady)="onPagingControllerReady($event)"
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
  pagingController: PagingController;
  totalSizeLabel = "unknown";

  isEnterPressed = isEnterKey;
  guessFileThumbnail = guessFileThumbnail;
  preview = (u) => {
    preview(u, this.dialog, this.nav, this.fileService, isMobile(), () => {
      this.dialogRef.close();
    });
  };

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    public dialogRef: MatDialogRef<
      VerFileHistoryComponent,
      VerFileHistoryDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: VerFileHistoryDialogData,
    private fileService: FileInfoService,
    private dialog: MatDialog,
    private nav: NavigationService
  ) {}

  ngOnInit(): void {
    this.qryTotalSize();
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.setPageLimit(5);
    this.pagingController.PAGE_LIMIT_OPTIONS = [5];
    this.pagingController.onPageChanged = () => this.fetch();
    this.fetch();
  }

  qryTotalSize() {
    this.http
      .post<any>(
        `${environment.vfm}/open/api/versioned-file/accumulated-size`,
        { verFileId: this.data.verFileId }
      )
      .subscribe({
        next: (r) => {
          this.totalSizeLabel = resolveSize(r.data.sizeInBytes);
        },
      });
  }

  fetch() {
    this.http
      .post<any>(`${environment.vfm}/open/api/versioned-file/history`, {
        paging: this.pagingController.paging,
        verFileId: this.data.verFileId,
      })
      .subscribe({
        next: (r) => {
          if (!r.data.payload) {
            r.data.payload = [];
          }
          this.tabdata = r.data.payload;
          this.pagingController.onTotalChanged(r.data.paging);
          for (let f of this.tabdata) {
            if (f.thumbnail) {
              f.thumbnail =
                environment.fstore +
                "/file/raw?key=" +
                encodeURIComponent(f.thumbnail);
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
          class="form-control"
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

    <div class="d-grid gap-2 d-md-flex justify-content-md-end mb-3">
      <button
        mat-raised-button
        class="m-2"
        [class.status-green]="expandUploadPanel"
        (click)="toggleUploadPanel()"
      >
        Upload Panal
      </button>
      <button mat-raised-button class="m-2" (click)="fetch()">Fetch</button>
      <button mat-raised-button class="m-2" (click)="reset()">Reset</button>
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
          <td mat-cell *matCellDef="let u" (click)="$event.stopPropagation()">
            {{ u.verFileId }}
          </td>
        </ng-container>

        <ng-container matColumnDef="thumbnail">
          <th mat-header-cell *matHeaderCellDef><b>Preview</b></th>
          <td mat-cell *matCellDef="let f" (click)="preview(f)">
            <img
              style="max-height:50px; padding: 5px 0px 5px 0px;"
              *ngIf="f.thumbnail"
              [src]="f.thumbnail"
            />
            <img
              style="max-height:40px; padding: 5px 0px 5px 0px;"
              *ngIf="!f.thumbnail"
              [src]="guessFileThumbnail(f)"
            />
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
            <button mat-raised-button class="m-2" (click)="selectVerFile(f)">
              Update
            </button>
            <button
              mat-raised-button
              class="m-2"
              (click)="showVerFileHistory(f)"
            >
              History
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="tabcol"></tr>
        <tr mat-row *matRowDef="let row; columns: tabcol"></tr>
      </table>
    </div>

    <app-controlled-paginator
      (controllerReady)="onPagingControllerReady($event)"
    ></app-controlled-paginator>
  `,
  styles: [],
})
export class VersionedFileComponent implements OnInit {
  expandUploadPanel = false;
  pagingController: PagingController;
  tabdat: ApiListVerFileRes[] = [];
  tabcol = ["thumbnail", "verFileId", "name", "uploadTime", "size", "operate"];
  searchName = "";
  updateVerFileId = "";
  updateVerFileName = "";

  isMobile = false;
  progress: string = null;
  isUploading: boolean = false;
  uploadFileName: string = null;
  uploadFile: File = null;
  uploadSub: Subscription = null;

  guessFileThumbnail = guessFileThumbnail;
  preview = (u) => {
    preview(u, this.dialog, this.nav, this.fileService, this.isMobile);
  };

  constructor(
    private http: HttpClient,
    private fileService: FileInfoService,
    private toaster: Toaster,
    private dialog: MatDialog,
    private nav: NavigationService
  ) {}

  ngOnInit(): void {
    this.isMobile = isMobile();
  }

  fetch() {
    let req: ApiListVerFileReq = {
      paging: this.pagingController.paging,
      name: this.searchName,
    };
    this.http
      .post<any>(`${environment.vfm}/open/api/versioned-file/list`, req)
      .subscribe({
        next: (r) => {
          if (!r.data.payload) {
            r.data.payload = [];
          }
          this.tabdat = r.data.payload;
          this.pagingController.onTotalChanged(r.data.paging);
          for (let f of this.tabdat) {
            if (f.thumbnail) {
              f.thumbnail =
                environment.fstore +
                "/file/raw?key=" +
                encodeURIComponent(f.thumbnail);
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetch();
    this.fetch();
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
    this.toaster.toast("File uploading cancelled");
  }

  upload() {
    if (!this.uploadFile) {
      this.toaster.toast("Please select file to upload");
      return;
    }

    const abortUpload = () => {
      this.progress = null;
      this.isUploading = false;
      this.uploadFile = null;
      this.uploadFileName = "";
      this.toaster.toast(`Failed to upload file`);
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
              sub = this.http.post(
                `${environment.vfm}/open/api/versioned-file/update`,
                {
                  filename: this.uploadFileName,
                  fstoreFileId: fstoreRes.data,
                  verFileId: this.updateVerFileId,
                }
              );
            } else {
              sub = this.http.post(
                `${environment.vfm}/open/api/versioned-file/create`,
                {
                  filename: this.uploadFileName,
                  fstoreFileId: fstoreRes.data,
                }
              );
            }

            sub.subscribe({
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
}
