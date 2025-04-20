import { NestedTreeControl } from "@angular/cdk/tree";
import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { DirTopDownTreeNode, DirTree } from "src/common/dir-tree";
import { Paging, PagingController } from "src/common/paging";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { Env } from "src/common/env-util";

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
  attempt?: number;
  url?: string;
  platform?: string;
  dirFileKey?: string;
  dirName?: string;
  createdAt?: number;
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
  dirTreeControl = new NestedTreeControl<DirTopDownTreeNode>(
    (node) => node.child
  );
  platforms: string[] = [];
  dirTreeDataSource = new MatTreeNestedDataSource<DirTopDownTreeNode>();
  pagingController: PagingController;
  tabdata: ListedTask[] = [];

  ngOnInit(): void {
    this.headers = this.env.isMobile()
      ? ["status", "url", "operation"]
      : [
          "taskId",
          "status",
          "url",
          "platform",
          "attempt",
          "dirName",
          "fileCount",
          "createdAt",
          "remark",
          "operation",
        ];
  }

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    public dirTree: DirTree,
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
        if (this.tabdata) {
          for (let t of this.tabdata) {
            if (t.remark && t.remark.length > 100) {
              t.remark = "..." + t.remark.substring(t.remark.length - 100);
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.listTasks();
    this.listTasks();
  }

  fetchTopDownDirTree() {
    this.dirTree.fetchTopDownDirTree((dat) => {
      this.dirTreeDataSource.data = [dat];
      this.dirTreeControl.dataNodes = this.dirTreeDataSource.data;

      this.http.get<any>(`/drone/open/api/task/last-selected-dir`).subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          let dat: FetchLastSelectedDirRes = resp.data;
          if (dat.fileKey) {
            let children = this.dirTreeDataSource.data;
            while (children.length > 0) {
              let c = children.shift();
              if (c.fileKey && c.fileKey == dat.fileKey) {
                this.selectDir(c);
                break;
              }
              if (c.child) {
                children.push(...c.child);
              }
            }
          }
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
    });
  }

  selectDir(n) {
    this.createTaskReq.dirFileKey = n.fileKey;
    this.createTaskDirName = n.name;
    this.dirTreeControl.collapseAll();
  }

  toggleCreateTaskPanelShown() {
    this.createTaskPanelShown = !this.createTaskPanelShown;
    if (this.createTaskPanelShown) {
      this.fetchTopDownDirTree();
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
      this.guessUrlPlatform();
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
}
