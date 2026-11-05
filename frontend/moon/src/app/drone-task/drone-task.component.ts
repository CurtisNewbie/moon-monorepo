import { HttpClient } from "@angular/common/http";
import { ActivatedRoute } from "@angular/router";
import {
  Component,
  ElementRef,
  Inject,
  OnDestroy,
  OnInit,
  ViewChild,
} from "@angular/core";
import { MatMenuTrigger } from "@angular/material/menu";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Paging } from "src/common/paging";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { isEnterKey } from "src/common/condition";
import { copyToClipboard } from "src/common/clipboard";
import { dialogSize } from "src/common/dialog-size";
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
  dirName?: string;
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
  makeDirName?: string;
  createdAt?: number;
  updatedAt?: number;
  fileCount?: number;
  lastPage?: number;
  remark?: string;
  remarkShort?: string;
  thumbnailFstToken?: string;
  thumbnail?: string;
}

export interface UpdateTaskReq {
  taskId?: string; // Required.
  url?: string;
  lastPage?: number;
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
export class DroneTaskComponent implements OnInit, OnDestroy {
  private readonly compactColumns = ["thumbnail", "status", "dirName"];
  private readonly wideColumns = [
    "thumbnail",
    "taskId",
    "status",
    "url",
    "platform",
    "attempt",
    "lastPage",
    "dirName",
    "fileCount",
    "updatedAt",
    "remark",
    "operation",
  ];
  displayedColumns: string[] = this.wideColumns;
  private readonly modeSub = this.env.mobile$.subscribe((mobile) => {
    this.displayedColumns = mobile ? this.compactColumns : this.wideColumns;
  });

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
  selectedFileKey: string | undefined;
  platforms: string[] = [];
  lastPagePlatforms: string[] = [];
  tabdata: ListedTask[] = [];
  listTaskReq: ListTaskReq = {};
  isEnterKey = isEnterKey;

  /** hidden anchor carrying the row menu; it opens where the touch lands */
  @ViewChild("rowMenuAnchor", { static: true })
  rowMenuAnchor: ElementRef;
  @ViewChild(MatMenuTrigger, { static: true })
  rowMenuTrigger: MatMenuTrigger;
  /** viewport position of the hidden menu anchor */
  rowMenuPosition = { x: "0", y: "0" };
  /** the row the open menu acts on */
  menuRow: ListedTask = null;
  /** pending long-press timer, cleared when the press is cancelled */
  private longPressTimer: any = null;
  private longPressStart = { x: 0, y: 0 };

  @ViewChild(ControlledPaginatorComponent, { static: true })
  pagingController: ControlledPaginatorComponent;

  ngOnInit(): void {
    this.pagingController.setPageLimitOptions([5, 10, 30, 50, 100]);
    this.listPlatforms();
    this.listLastPagePlatforms();
    const prefillUrl = this.route.snapshot.queryParams['url'] || '';
    if (prefillUrl) {
      this.bulkUrlFields = [prefillUrl];
      this.createTaskPanelShown = true;
    }
  }

  ngOnDestroy(): void {
    this.modeSub.unsubscribe();
  }

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private nav: NavigationService,
    public env: Env,
    private dialog: MatDialog,
    private i18n: I18n,
    private route: ActivatedRoute
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

          const remarkShortMaxLen = 40;
          const remarkMaxLen = 800;

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
          this.selectedFileKey = dat?.fileKey || undefined;
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
      this.selectedFileKey = undefined;
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

  listLastPagePlatforms() {
    this.http.get<any>(`/drone/open/api/list-last-page-platforms`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.lastPagePlatforms = resp.data || [];
      },
      error: (err) => {
        console.log(err);
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

  updateTask(req: UpdateTaskReq) {
    this.http.post<any>(`/drone/open/api/task/update`, req).subscribe({
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

  /**
   * Which actions a row offers. The wide buttons and the compact row menu read
   * the same predicates, so a phone never shows an action the row lacks.
   */
  canCancel(u: ListedTask): boolean {
    return u.status == "PENDING" || u.status == "UNRESOLVED";
  }

  canResolve(u: ListedTask): boolean {
    return u.status == "UNRESOLVED" || u.status == "RESOLVE_FAILED";
  }

  canRetry(u: ListedTask): boolean {
    return (
      u.status != "PENDING" &&
      u.status != "UNRESOLVED" &&
      u.status != "RESOLVE_FAILED" &&
      !!u.dirFileKey
    );
  }

  canUpdate(u: ListedTask): boolean {
    return (u.status == "COMPLETED" || u.status == "CANCELLED") && !!u.dirFileKey;
  }

  /** Whether the row's menu has anything to say; an empty menu is not opened. */
  canShowRowMenu(u: ListedTask): boolean {
    return (
      !!u &&
      (!!u.dirFileKey ||
        this.canCancel(u) ||
        this.canResolve(u) ||
        this.canRetry(u) ||
        this.canUpdate(u))
    );
  }

  /** Mobile entry: press and hold a row for half a second */
  onRowTouchStart(event: TouchEvent, row: ListedTask): void {
    if (!this.env.isMobile() || !this.canShowRowMenu(row)) {
      return;
    }
    const touch = event.touches[0];
    if (!touch) {
      return;
    }
    // a press that starts on a button is that button's, not the row's
    const target = event.target as Element;
    if (target?.closest?.("button")) {
      return;
    }
    this.clearLongPress();
    this.longPressStart = { x: touch.clientX, y: touch.clientY };
    this.longPressTimer = setTimeout(() => {
      this.longPressTimer = null;
      this.openRowMenu(row, touch.clientX, touch.clientY);
    }, 500);
  }

  /** A finger that travels is scrolling, not holding */
  onRowTouchMove(event: TouchEvent): void {
    if (!this.longPressTimer) {
      return;
    }
    const touch = event.touches[0];
    if (!touch) {
      return;
    }
    if (
      Math.abs(touch.clientX - this.longPressStart.x) > 10 ||
      Math.abs(touch.clientY - this.longPressStart.y) > 10
    ) {
      this.clearLongPress();
    }
  }

  onRowTouchEnd(): void {
    this.clearLongPress();
  }

  private openRowMenu(row: ListedTask, x: number, y: number): void {
    if (this.rowMenuTrigger.menuOpen) {
      return;
    }
    this.menuRow = row;
    this.rowMenuPosition = { x: x + "px", y: y + "px" };
    // the overlay measures the anchor as it attaches, and change detection has
    // not flushed the position binding yet, so write it to the element first
    const anchor = this.rowMenuAnchor?.nativeElement as HTMLElement;
    if (anchor) {
      anchor.style.left = this.rowMenuPosition.x;
      anchor.style.top = this.rowMenuPosition.y;
    }
    this.rowMenuTrigger.openMenu();
  }

  private clearLongPress(): void {
    if (this.longPressTimer) {
      clearTimeout(this.longPressTimer);
      this.longPressTimer = null;
    }
  }

  clickUpdateTaskUrl(d) {
    const lastPageSupported = this.lastPagePlatforms.includes(d.platform);
    this.dialog
      .open(UpdateDroneTaskDialogComponent, {
        ...dialogSize(700),
        data: {
          taskId: d.taskId,
          url: d.url,
          lastPage: d.lastPage,
          lastPageSupported: lastPageSupported,
        },
      })
      .afterClosed()
      .subscribe((result) => {
        if (!result) {
          return;
        }
        let req: UpdateTaskReq = { taskId: d.taskId };
        if (result.url && result.url != d.url) {
          req.url = result.url;
        }
        if (
          lastPageSupported &&
          result.lastPage &&
          result.lastPage >= 1 &&
          result.lastPage != d.lastPage
        ) {
          req.lastPage = result.lastPage;
        }
        if (req.url === undefined && req.lastPage === undefined) {
          return;
        }
        this.updateTask(req);
      });
  }

  resolveTask(d) {
    this.dialog
      .open(ResolveDroneTaskDialogComponent, {
        ...dialogSize(600),
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
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'taskId' }}</mat-label>
        <input readonly disabled matInput [(ngModel)]="data.taskId" />
      </mat-form-field>
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'url' }}</mat-label>
        <input matInput [(ngModel)]="data.url" />
      </mat-form-field>
      <mat-form-field style="width: 100%; max-width: 400px" *ngIf="data.lastPageSupported">
        <mat-label>{{ 'drone-task' | trl:'lastPage' }}</mat-label>
        <input matInput type="number" min="1" [(ngModel)]="data.lastPage" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="{url: data.url, lastPage: data.lastPage}">
        {{ 'drone-task' | trl:'update' }}
      </button>
      <button mat-button [mat-dialog-close]="null" cdkFocusInitial>
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
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'taskId' }}</mat-label>
        <input readonly disabled matInput [ngModel]="data.taskId" />
      </mat-form-field>
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'url' }}</mat-label>
        <input readonly disabled matInput [ngModel]="data.url" />
      </mat-form-field>
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'platform' }}</mat-label>
        <mat-select [(ngModel)]="data.platform">
          <mat-option [value]="option" *ngFor="let option of data.platforms">
            {{option}}
          </mat-option>
        </mat-select>
      </mat-form-field>
      <mat-form-field style="width: 100%; max-width: 400px">
        <mat-label>{{ 'drone-task' | trl:'newDirectoryNameOptional' }}</mat-label>
        <input matInput [(ngModel)]="data.makeDirName" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="{platform: data.platform, makeDirName: data.makeDirName}">
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
