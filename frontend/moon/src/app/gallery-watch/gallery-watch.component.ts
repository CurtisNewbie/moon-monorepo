import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Env } from "src/common/env-util";
import { I18n } from "../i18n.service";
import { NavigationService } from "../navigation.service";

interface WatchTaskReq {
  keywords: string;
  platform: string;
}
interface WatchTaskRes {
  taskId: string;
}
interface ListedWatchTask {
  taskId: string;
  keywords: string;
  platform: string;
  enabled: boolean;
  maxItems?: number;
}
interface EnableWatchTaskReq {
  taskId: string;
  enabled: boolean;
}
interface TriggerWatchTaskReq {
  taskId: string;
}

@Component({
  selector: "app-gallery-watch",
  templateUrl: "./gallery-watch.component.html",
  styleUrls: ["./gallery-watch.component.css"],
})
export class GalleryWatchComponent implements OnInit {
  headers: string[] = ["platform", "keywords", "maxItems", "enabled", "operation"];

  // Create Watch Task
  selectedPlatform: string = "";
  keywords: string = "";
  platforms: string[] = [];
  showCreatePanel: boolean = false;

  // Watch Task list
  watchTasks: ListedWatchTask[] = [];

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    public env: Env,
    private i18n: I18n,
    private nav: NavigationService
  ) {}

  ngOnInit(): void {
    this.listPlatforms();
    this.listWatchTasks();
  }

  listPlatforms() {
    this.http.get<any>(`/drone/open/api/list-platforms`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.platforms = resp.data;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  listWatchTasks() {
    this.http.get<any>(`/drone/open/api/list-watch-tasks`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.watchTasks = resp.data || [];
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  createWatchTask() {
    if (!this.keywords.trim() || !this.selectedPlatform) {
      this.snackBar.open("Keywords and platform are required", "ok", {
        duration: 3000,
      });
      return;
    }
    const req: WatchTaskReq = {
      keywords: this.keywords.trim(),
      platform: this.selectedPlatform,
    };
    this.http.post<any>(`/drone/open/api/create-watch-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        const dat: WatchTaskRes = resp.data;
        this.snackBar.open(`Watch task created: ${dat.taskId}`, "ok", {
          duration: 5000,
        });
        this.keywords = "";
        this.selectedPlatform = "";
        this.showCreatePanel = false;
        this.listWatchTasks();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  toggleCreatePanel(): void {
    this.showCreatePanel = !this.showCreatePanel;
  }

  viewGalleries(taskId: string) {
    this.nav.navigateToWithPreserve('/gallery-watch/galleries', { taskId: taskId });
  }

  toggleTaskEnabled(task: ListedWatchTask): void {
    const newEnabled = !task.enabled;
    const req: EnableWatchTaskReq = {
      taskId: task.taskId,
      enabled: newEnabled,
    };
    this.http.post<any>('/drone/open/api/toggle-watch-task', req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
          return;
        }
        task.enabled = newEnabled;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open('Request failed, unknown error', 'ok', {
          duration: 3000,
        });
      },
    });
  }

  triggerWatchTask(taskId: string): void {
    const req: TriggerWatchTaskReq = { taskId: taskId };
    this.http.post<any>(`/drone/open/api/trigger-watch-task`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.snackBar.open("Triggered", "ok", { duration: 3000 });
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
