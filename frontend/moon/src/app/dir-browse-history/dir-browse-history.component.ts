import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";

export interface DirBrowseRecord {
  dirKey: string;
  name: string;
  thumbnailToken?: string;
  fileKey: string;
  time: number; // unix millis
}

@Component({
  selector: "app-dir-browse-history",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Directory Browse History</h3>
    </div>
    <mat-divider></mat-divider>
    <div class="container m-4">
      <cdk-virtual-scroll-viewport itemSize="286" style="height: 85vh">
        <div *cdkVirtualFor="let it of dat">
          <mat-card class="mat-elevation-z2 m-3">
            <div>
              <div class="row">
                <div class="col-md-4">
                  <div class="m-1">
                    <img
                      *ngIf="it.thumbnailToken"
                      style="height:200px"
                      class="m-2 mat-elevation-z8 p-1"
                      [src]="thumbnailUrl(it)"
                    />
                    <i
                      *ngIf="!it.thumbnailToken"
                      class="bi bi-folder icon-button-large-preview mat-elevation-z8"
                    ></i>
                  </div>
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
                    <mat-label>Last File</mat-label>
                    <input
                      matInput
                      type="text"
                      [ngModel]="it.fileKey"
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
                    />
                  </mat-form-field>
                  <div class="m-2" matLine>
                    <button mat-icon-button (click)="resumeReading(it.dirKey, it.fileKey)">
                      Resume <i class="bi bi-arrow-right-circle"></i>
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
export class DirBrowseHistoryComponent implements OnInit {
  dat: DirBrowseRecord[] = [];

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private navigation: NavigationService
  ) {}

  ngOnInit(): void {
    this.fetchHistory();
  }

  fetchHistory() {
    this.http.get<any>(`/vfm/history/list-dir-browse`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: DirBrowseRecord[] = resp.data;
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

  resumeReading(dirKey: string, fileKey: string) {
    this.navigation.navigateTo(NavType.MANAGE_FILES, [
      { parentDirKey: dirKey, targetFileKey: fileKey, orderBy: 'name', autoPreview: true },
    ]);
  }

  thumbnailUrl(it: DirBrowseRecord) {
    return "fstore/file/raw?key=" + encodeURIComponent(it.thumbnailToken!);
  }
}
