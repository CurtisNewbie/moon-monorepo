import { HttpClient } from "@angular/common/http";
import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Paging } from "src/common/paging";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { DirTreeNavComponent } from "../dir-tree-nav/dir-tree-nav.component";

export interface ApiDetectTitleReq {
  url?: string;
  platform?: string;
}

export interface ApiDetectTitleRes {
  title?: string;
  platform?: string;
}

export interface GuessUrlPlatformRes {
  platform?: string;
}

export interface GuessUrlPlatformReq {
  url?: string;
}

export interface FetchLastSelectedDirRes {
  fileKey?: string;
}

export interface RetryTaskReq {
  taskId?: string;
}

export interface CreateTaskReq {
  dirFileKey?: string;
  platform?: string;
  url?: string;
  makeDirName?: string;
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
}

@Component({
  selector: "app-drone-task",
  templateUrl: "./drone-task.component.html",
  styleUrls: ["./drone-task.component.css"],
})
export class DroneTaskComponent implements OnInit {
  headers: string[] = [];
  createTaskPanelShown: boolean = false;
  createTaskDirName: string = "";
  createTaskReq: CreateTaskReq = {};
  platforms: string[] = [];
  tabdata: ListedTask[] = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  @ViewChild(DirTreeNavComponent)
  dirTreeNav: DirTreeNavComponent;

  ngOnInit(): void {
    this.headers = this.env.isMobile()
      ? ["status", "platform", "dirName", "operation"]
      : [
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
  }

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private nav: NavigationService,
    public env: Env
  ) {}

  listTasks() {
    let req: ListTaskReq = { paging: this.pagingController.paging };
    this.http.post<any>(`/drone/open/api/list-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: PageRes = resp.data;
        this.tabdata = dat.payload;

        const isMob = this.env.isMobile();
        const remarkMaxLen = 40;
        const statusLabelMaxLen = 60;

        if (this.tabdata) {
          for (let t of this.tabdata) {
            if (t.remark && t.remark.length > remarkMaxLen) {
              t.remark =
                "... " +
                t.remark.substring(t.remark.length - remarkMaxLen).trim();
            }
            t.trimmedDirName = t.dirName;
            if (isMob && t.trimmedDirName.length > statusLabelMaxLen) {
              t.trimmedDirName =
                t.trimmedDirName.substring(0, statusLabelMaxLen) + " ...";
            }
            t.statusLabel = t.status;
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
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  createTask() {
    let req: CreateTaskReq = this.createTaskReq;
    this.http.post<any>(`/drone/open/api/create-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.createTaskReq = {};
        this.createTaskDirName = "";
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

  toggleCreateTaskPanelShown() {
    this.createTaskPanelShown = !this.createTaskPanelShown;
    if (this.createTaskPanelShown) {
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

  urlTypingTimer = null;

  urlKeyUp() {
    if (this.urlTypingTimer != null) {
      window.clearTimeout(this.urlTypingTimer);
    }
    this.urlTypingTimer = window.setTimeout(() => {
      this.createTaskReq.url = this.createTaskReq.url.trim();
      if (!this.createTaskReq.platform && this.createTaskReq.makeDirName) {
        this.guessUrlPlatform();
      } else {
        this.detectTitle();
      }
    }, 500);
  }

  urlKeyDown() {
    if (this.urlTypingTimer != null) {
      window.clearTimeout(this.urlTypingTimer);
    }
  }

  guessUrlPlatform() {
    if (!this.createTaskReq.url) {
      return;
    }
    let req: GuessUrlPlatformReq = { url: this.createTaskReq.url };
    this.http.post<any>(`/drone/open/api/task/guess-plafrom`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: GuessUrlPlatformRes = resp.data;
        if (dat && dat.platform) {
          this.createTaskReq.platform = dat.platform;
        }
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  detectTitle() {
    if (!this.createTaskReq.url || this.createTaskReq.makeDirName) {
      return;
    }

    let req: ApiDetectTitleReq = {
      url: this.createTaskReq.url,
      platform: this.createTaskReq.platform,
    };
    this.http.post<any>(`/drone/open/api/task/detect-title`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: ApiDetectTitleRes = resp.data;
        if (!this.createTaskReq.makeDirName) {
          this.createTaskReq.makeDirName = dat.title;
        }
        if (!this.createTaskReq.platform) {
          this.createTaskReq.platform = dat.platform;
        }
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }
}
