import { HttpClient } from "@angular/common/http";
import { Component, OnInit, ViewChild } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { Env } from "src/common/env-util";
import { I18n } from "../i18n.service";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { NavigationService } from "../navigation.service";

interface ListGalleryReq {
  paging: { page: number; limit: number };
  taskId: string;
}
interface ListedGallery {
  url: string;
  name: string;
  thumbnailToken: string;
  createdAt: string;
  thumbnailUrl?: string;
  taskCreated?: boolean;
}
interface ListGalleryRes {
  data: {
    paging: { page: number; limit: number; total: number };
    payload: ListedGallery[];
  };
}

@Component({
  selector: "app-watched-gallery",
  templateUrl: "./watched-gallery.component.html",
  styleUrls: ["./watched-gallery.component.css"],
})
export class WatchedGalleryComponent implements OnInit {
  headers: string[] = ["thumbnail", "name", "url", "createdAt", "operation"];

  taskId: string = "";
  platform: string = '';
  keywords: string = '';
  galleries: ListedGallery[] = [];
  listGalleryReq: ListGalleryReq = {
    paging: { page: 1, limit: 10 },
    taskId: "",
  };
  total: number = 0;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    public env: Env,
    private i18n: I18n,
    private nav: NavigationService,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.taskId = this.route.snapshot.queryParams['taskId'] || '';
    this.listGalleryReq.taskId = this.taskId;
    this.fetchWatchTask();
    this.searchGalleries();
  }

  searchGalleries() {
    this.listGalleryReq.paging = this.pagingController
      ? this.pagingController.paging
      : { page: 1, limit: 10 };
    this.listGalleryReq.taskId = this.taskId;
    this.http
      .post<any>(`/drone/open/api/list-watched-gallery`, this.listGalleryReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          const dat = resp.data;
          this.galleries = dat.payload || [];
          this.total = dat.paging?.total || 0;
          if (this.galleries) {
            for (let g of this.galleries) {
              if (g.thumbnailToken) {
                g.thumbnailUrl =
                  "fstore/file/raw?key=" +
                  encodeURIComponent(g.thumbnailToken);
              }
            }
          }
          if (this.pagingController) {
            this.pagingController.onTotalChanged(dat.paging);
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

  onPageChanged() {
    this.searchGalleries();
  }

  fetchWatchTask(): void {
    if (!this.taskId) return;
    this.http.get<any>('/drone/open/api/list-watch-tasks').subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
          return;
        }
        const tasks = resp.data || [];
        const task = tasks.find((t: any) => t.taskId === this.taskId);
        if (task) {
          this.platform = task.platform || '';
          this.keywords = task.keywords || '';
        }
      },
      error: (err: any) => {
        console.log(err);
        this.snackBar.open('Request failed, unknown error', 'ok', {
          duration: 3000,
        });
      },
    });
  }

  goToDroneTask(url: string): void {
    this.nav.navigateToWithPreserve('/drone-task', { url: url });
  }

  goBack() {
    this.nav.navigateToUrl('/gallery-watch');
  }
}
