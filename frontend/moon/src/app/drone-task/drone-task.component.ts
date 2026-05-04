import { HttpClient } from "@angular/common/http";
import { Component, Inject, OnInit, AfterViewInit, ViewChild } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Paging } from "src/common/paging";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { DirTreeNavComponent } from "../dir-tree-nav/dir-tree-nav.component";
import { isEnterKey } from "src/common/condition";
import { copyToClipboard } from "src/common/clipboard";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { I18n } from "../i18n.service";

export interface CreateTaskReq {
  dirFileKey?: string;
  platform?: string;
  url?: string;
  makeDirName?: string;
}

export interface CreateTasksBatchReq {
  dirFileKey?: string;
  platform?: string;
  urls?: string[];
}

export interface RetryTaskReq {
  taskId?: string;
}

export interface CancelTaskReq {
  taskId?: string;
}

export interface PageRes {
  paging?: Paging;
  payload?: ListedTask[];
}

export interface ListTaskReq {
  paging?: Paging;
  url?: string;
  taskId?: string;
  status?: string;
  platform?: string;
}
export interface ListedTask {
  taskId?: string;
  status?: string;
  statusLabel?: string;
  attempt?: number;
  url?: string;
  platform?: string;
  dirFileKey?: string;
  dirName?: string;
  trimmedDirName?: string;
  createdAt?: number;
  updatedAt?: number;
  fileCount?: number;
  remark?: string;
  remarkShort?: string;
  thumbnailFstToken?: string;
  thumbnail?: string;
}

export interface UpdateTaskURLReq {
  taskId?: string; // Required.
  url?: string; // Required.
}

export interface ResolveTaskReq {
  taskId?: string;
  platform?: string;
  makeDirName?: string;
}

@Component({
  selector: "app-drone-task",
  templateUrl: "./drone-task.component.html",
  styleUrls: ["./drone-task.component.css"],
})
export class DroneTaskComponent implements OnInit, AfterViewInit {
  headers: string[] = [];
  createTaskPanelShown: boolean = false;
  bulkUrlFields: string[] = [""];

  get parsedUrls(): string[] {
    const seen = new Set<string>();
    return this.bulkUrlFields
      .map((u) => u.trim())
      .filter((u) => u.length > 0 && !seen.has(u) && seen.add(u));
  }

  addUrlField() {
    this.bulkUrlFields.push("");
  }

  removeUrlField(i: number) {
    if (this.bulkUrlFields.length > 1) {
      this.bulkUrlFields.splice(i, 1);
    }
  }
  createTaskDirName: string = "";
  createTaskReq: CreateTaskReq = {};
  prevSelectedFileKey: string | undefined;
  platforms: string[] = [];
  tabdata: ListedTask[] = [];
  listTaskReq: ListTaskReq = {};
  isEnterKey = isEnterKey;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  @ViewChild(DirTreeNavComponent)
  dirTreeNav: DirTreeNavComponent;

  ngOnInit(): void {
    this.headers = this.env.isMobile()
      ? ["thumbnail", "status", "dirName", "operation"]
      : [
          "thumbnail",
          "taskId",
          "status",
          "url",
          "platform",
          "attempt",
          "dirName",
          "fileCount",
          "updatedAt",
          "remark",
          "operation",
        ];
    this.listPlatforms();
  }

  ngAfterViewInit(): void {
    this.pagingController.setPageLimitOptions([5, 10, 30, 50, 100]);
  }

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private nav: NavigationService,
    public env: Env,
    private dialog: MatDialog,
    private i18n: I18n
  ) {}

  listTasks() {
    this.listTaskReq.paging = this.pagingController.paging;
    this.http
      .post<any>(`/drone/open/api/list-task`, this.listTaskReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          let dat: PageRes = resp.data;
          this.tabdata = dat.payload;

          const isMob = this.env.isMobile();
          const remarkShortMaxLen = 40;
          const remarkMaxLen = 800;
          const statusLabelMaxLen = 60;

          if (this.tabdata) {
            for (let t of this.tabdata) {
              t.remarkShort = t.remark;
              if (t.remark && t.remark.length > remarkShortMaxLen) {
                t.remarkShort =
                  "... " +
                  t.remark
                    .substring(t.remark.length - remarkShortMaxLen)
                    .trim();
                if (t.remark.length > remarkMaxLen) {
                  t.remark =
                    "... " +
                    t.remark.substring(t.remark.length - remarkMaxLen).trim();
                }
              }
              if (t.thumbnailFstToken) {
                t.thumbnail =
                  "fstore/file/raw?key=" +
                  encodeURIComponent(t.thumbnailFstToken);
              }
              t.trimmedDirName = t.dirName;
              if (isMob && t.trimmedDirName.length > statusLabelMaxLen) {
                t.trimmedDirName =
                  t.trimmedDirName.substring(0, statusLabelMaxLen) + " ...";
              }
              t.statusLabel = t.status;
              if (t.status == 'COMPLETED') {
                t.statusLabel = this.i18n.trl('drone-task', 'completed');
              } else if (t.status == 'PENDING') {
                t.statusLabel = this.i18n.trl('drone-task', 'pending');
              } else if (t.status == 'UNRESOLVED') {
                t.statusLabel = this.i18n.trl('drone-task', 'unresolved');
              } else if (t.status == 'RESOLVE_FAILED') {
                t.statusLabel = this.i18n.trl('drone-task', 'resolveFailed');
              } else if (t.status == 'CANCELLED') {
                t.statusLabel = this.i18n.trl('drone-task', 'cancelled');
              }
            }
          }
          this.pagingController.onTotalChanged(dat.paging);
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  cancelTask(taskId) {
    let req: CancelTaskReq = { taskId: taskId };
    this.http.post<any>(`/drone/open/api/cancel-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }

        this.listTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  trackByIndex(i: number) {
    return i;
  }

  singleSubmit() {
    this.createTaskReq.url = (this.bulkUrlFields[0] || "").trim();
    if (!this.createTaskReq.url) {
      return;
    }
    let req: CreateTaskReq = this.createTaskReq;
    this.http.post<any>(`/drone/open/api/create-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.createTaskReq = {};
        this.createTaskDirName = "";
        this.bulkUrlFields = [""];
        this.createTaskPanelShown = false;
        this.listTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
      });
  }

  createTasksBatch() {
    const urls = this.parsedUrls;
    if (urls.length === 0) {
      return;
    }
    let req: CreateTasksBatchReq = {
      dirFileKey: this.createTaskReq.dirFileKey,
      platform: this.createTaskReq.platform,
      urls: urls,
    };
    this.http.post<any>(`/drone/open/api/create-tasks-batch`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.snackBar.open(
          `Created ${resp.data} tasks`,
          "ok",
          { duration: 3000 }
        );
        this.createTaskReq = {};
        this.createTaskDirName = "";
        this.bulkUrlFields = [""];
        this.createTaskPanelShown = false;
        this.listTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  selectDir(n) {
    this.createTaskReq.dirFileKey = n.fileKey;
    this.createTaskDirName = n.name;
    this.dirTreeNav.collapseAll();
  }

  fetchLastSelectedDir() {
    this.http
      .get<any>(`/drone/open/api/task/last-selected-dir`)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          let dat: { fileKey?: string } = resp.data;
          this.prevSelectedFileKey = dat?.fileKey || undefined;
        },
        error: (err) => {
          console.log(err);
        },
      });
  }

  onTreeLoaded() {
    this.fetchLastSelectedDir();
  }

  toggleCreateTaskPanelShown() {
    this.createTaskPanelShown = !this.createTaskPanelShown;
    if (this.createTaskPanelShown) {
      this.prevSelectedFileKey = undefined;
      this.listPlatforms();
    }
  }

  listPlatforms() {
    this.http.get<any>(`/drone/open/api/list-platforms`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: string[] = resp.data;
        this.platforms = dat;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  gotoDir(dirFileKey) {
    this.nav.navigateTo(NavType.MANAGE_FILES, [{ parentDirKey: dirFileKey }]);
  }

  retryTask(taskId) {
    let req: RetryTaskReq = { taskId: taskId };
    this.http.post<any>(`/drone/open/api/retry-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.listTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  updateTaskUrl(req: UpdateTaskURLReq) {
    this.http.post<any>(`/drone/open/api/task/update-url`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.listTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  copyRemark(r) {
    copyToClipboard(r);
    this.snackBar.open("Copied to clipboard", "ok", { duration: 3000 });
  }

  clickUpdateTaskUrl(d) {
    this.dialog
      .open(UpdateDroneTaskDialogComponent, {
        width: "700px",
        data: {
          taskId: d.taskId,
          url: d.url,
        },
      })
      .afterClosed()
      .subscribe((newUrl) => {
        if (newUrl && newUrl != d.url) {
          this.updateTaskUrl({
            taskId: d.taskId,
            url: newUrl,
          });
        }
      });
  }

  resolveTask(d) {
    this.dialog
      .open(ResolveDroneTaskDialogComponent, {
        width: "600px",
        data: {
          taskId: d.taskId,
          url: d.url,
          platform: d.platform,
          makeDirName: d.makeDirName || "",
          platforms: this.platforms,
        },
      })
      .afterClosed()
      .subscribe((result) => {
        if (result) {
          let req: ResolveTaskReq = {
            taskId: d.taskId,
            platform: result.platform,
            makeDirName: result.makeDirName || undefined,
          };
          this.http
            .post<any>(`/drone/open/api/task/resolve`, req)
            .subscribe({
              next: (resp) => {
                if (resp.error) {
                  this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                  return;
                }
                this.listTasks();
              },
              error: (err) => {
                console.log(err);
                this.snackBar.open("Request failed, unknown error", "ok", {
                  duration: 3000,
                });
              },
            });
        }
      });
  }

  }

@Component({
  selector: "update-drone-task-component",
  template: `
    <h1 mat-dialog-title>{{ 'drone-task' | trl:'updateDroneTask' }}</h1>
    <div mat-dialog-content>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'taskId' }}</mat-label>
        <input readonly disabled matInput [(ngModel)]="data.taskId" />
      </mat-form-field>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'url' }}</mat-label>
        <input matInput [(ngModel)]="data.url" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="data.url">
        {{ 'drone-task' | trl:'update' }}
      </button>
      <button mat-button [mat-dialog-close]="''" cdkFocusInitial>
        {{ 'drone-task' | trl:'no' }}
      </button>
    </div>
  `,
})
export class UpdateDroneTaskDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<UpdateDroneTaskDialogComponent, any>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}
}

@Component({
  selector: "resolve-drone-task-component",
  template: `
    <h1 mat-dialog-title>{{ 'drone-task' | trl:'resolve' }}</h1>
    <div mat-dialog-content>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'taskId' }}</mat-label>
        <input readonly disabled matInput [ngModel]="data.taskId" />
      </mat-form-field>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'url' }}</mat-label>
        <input readonly disabled matInput [ngModel]="data.url" />
      </mat-form-field>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'platform' }}</mat-label>
        <mat-select [(ngModel)]="data.platform">
          <mat-option [value]="option" *ngFor="let option of data.platforms">
            {{option}}
          </mat-option>
        </mat-select>
      </mat-form-field>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'drone-task' | trl:'newDirectoryNameOptional' }}</mat-label>
        <input matInput [(ngModel)]="data.makeDirName" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="data.platform && data.makeDirName ? {platform: data.platform, makeDirName: data.makeDirName} : null">
        {{ 'drone-task' | trl:'resolve' }}
      </button>
      <button mat-button [mat-dialog-close]="null" cdkFocusInitial>
        {{ 'drone-task' | trl:'no' }}
      </button>
    </div>
  `,
})
export class ResolveDroneTaskDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<ResolveDroneTaskDialogComponent, any>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}
}
