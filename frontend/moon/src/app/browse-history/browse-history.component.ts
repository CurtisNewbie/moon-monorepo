import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { guessFileIconClz } from "src/common/file";

export interface ListBrowseRecordRes {
  time?: number;
  fileKey?: string;
  name?: string;
  thumbnailToken?: string;
  deleted?: boolean;
}

@Component({
  selector: "app-browse-history",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Browse History</h3>
    </div>
    <mat-divider></mat-divider>
    <div class="container m-4">
      <cdk-virtual-scroll-viewport itemSize="174" style="height: 80vh">
        <div *cdkVirtualFor="let it of dat">
          <mat-card class="mat-elevation-z2 m-3">
            <div>
              <div class="row">
                <div class="col">
                  <img
                    *ngIf="it.thumbnailToken"
                    style="height:120px"
                    class="m-2 mat-elevation-z8 p-3"
                    [src]="thumbnailUrl(it)"
                  />
                  <bi
                    *ngIf="!it.thumbnailToken"
                    [ngClass]="[
                      'mat-elevation-z8',
                      'icon-button-large-preview',
                      guessFileIcon(it)
                    ]"
                  ></bi>
                </div>
                <div class="col">
                  <mat-form-field style="width: 100%;" class="m-2">
                    <mat-label>Name</mat-label>
                    <input
                      matInput
                      type="text"
                      [ngModel]="it.name"
                      readonly="true"
                    />
                  </mat-form-field>
                  <mat-form-field style="width: 100%;" class="m-2">
                    <mat-label>Browse Time</mat-label>
                    <input
                      matInput
                      type="text"
                      [ngModel]="it.time | date : 'yyyy-MM-dd HH:mm:ss'"
                      readonly="true"
                      (ngModelChange)="it.value = $event"
                    />
                  </mat-form-field>
                  <div class="col" *ngIf="it.deleted">
                    <span class="status-red"><b>Deleted</b></span>
                  </div>
                  <div class="m-2" matLine *ngIf="!it.deleted">
                    <button mat-icon-button (click)="goToFile(it.fileKey)">
                      Find File <i class="bi bi-search"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </mat-card>
        </div>
      </cdk-virtual-scroll-viewport>
    </div>
    <mat-divider></mat-divider>
  `,
  styles: [],
})
export class BrowseHistoryComponent implements OnInit {
  guessFileIcon = guessFileIconClz;

  dat: ListBrowseRecordRes[] = [];

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private navigation: NavigationService
  ) {}

  ngOnInit(): void {
    this.fetchHistory();
  }

  fetchHistory() {
    this.http.get<any>(`/vfm/history/list-browse-history`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: ListBrowseRecordRes[] = resp.data;
        if (dat == null) {
          dat = [];
        }
        this.dat = dat;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  goToFile(fileKey) {
    this.navigation.navigateTo(NavType.MANAGE_FILES, [
      { searchedFileKey: fileKey },
    ]);
  }

  thumbnailUrl(it: ListBrowseRecordRes) {
    return "fstore/file/raw?key=" + encodeURIComponent(it.thumbnailToken);
  }
}
