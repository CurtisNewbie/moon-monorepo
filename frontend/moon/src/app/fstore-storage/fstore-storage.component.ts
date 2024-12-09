import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpClient } from "@angular/common/http";
import { UserService } from "../user.service";

export interface StorageUsageInfo {
  path?: string;
  used?: number;
  usedText?: string;
}

export interface StorageInfo {
  volumns?: VolumnInfo[];
}

export interface VolumnInfo {
  mounted?: string;
  total?: number;
  used?: number;
  available?: number;
  usedPercent?: number;
  totalText?: string;
  usedText?: string;
  availableText?: string;
  usedPercentText?: string;
}

@Component({
  selector: "app-fstore-storage",
  template: `
    <div>
      <h3 class="mt-2 mb-3">File Storage</h3>
    </div>

    <div class="d-flex flex-wrap justify-content-end gap-2">
      <ng-container *ngIf="hasRes('vfm:server:maintenance')">
        <button
          mat-raised-button
          (click)="compensateMissingThumbnails()"
          matTooltip="Generate missing thumbnails"
        >
          Compensate Missing Thumbnails
        </button>
        <button
          mat-raised-button
          (click)="regenerateVideoThumbnails()"
          matTooltip="Regenerate video thumbnails"
        >
          Regenerate Video Thumbnails
        </button>
      </ng-container>

      <ng-container *ngIf="hasRes('fstore:server:maintenance')">
        <button
          mat-raised-button
          (click)="sanitizeStorage()"
          matTooltip="Delete dangling files that are not uploaded to server"
        >
          Sanitize Storage
        </button>
        <button
          mat-raised-button
          (click)="removeDeletedFiles()"
          matTooltip="Prune files that are marked deleted. During server maintenance, file uploading is rejected"
        >
          Prune Deleted
        </button>
      </ng-container>
    </div>

    <div *ngIf="hasRes('fstore:fetch-storage-info')">
      <div class="d-flex justify-content-between">
        <h4 class="mt-3 mb-1">Storage Info</h4>
        <button mat-icon-button class="m-1 icon-button-large" (click)="fetch()">
          <i class="bi bi-arrow-clockwise"></i>
        </button>
      </div>

      <div class="mt-3 mb-2" style="overflow: auto;">
        <table mat-table [dataSource]="storageTabDat" style="width: 100%;">
          <ng-container matColumnDef="mounted">
            <th mat-header-cell *matHeaderCellDef>Mounted</th>
            <td mat-cell *matCellDef="let u">{{ u.mounted }}</td>
          </ng-container>

          <ng-container matColumnDef="total">
            <th mat-header-cell *matHeaderCellDef>Total</th>
            <td mat-cell *matCellDef="let u">{{ u.totalText }}</td>
          </ng-container>

          <ng-container matColumnDef="used">
            <th mat-header-cell *matHeaderCellDef>Used</th>
            <td mat-cell *matCellDef="let u">{{ u.usedText }}</td>
          </ng-container>

          <ng-container matColumnDef="available">
            <th mat-header-cell *matHeaderCellDef>Available</th>
            <td mat-cell *matCellDef="let u">{{ u.availableText }}</td>
          </ng-container>

          <ng-container matColumnDef="percent">
            <th mat-header-cell *matHeaderCellDef>Used Percentage</th>
            <td mat-cell *matCellDef="let u">{{ u.usedPercentText }}</td>
          </ng-container>

          <tr mat-header-row *matHeaderRowDef="storageTabCol"></tr>
          <tr mat-row *matRowDef="let row; columns: storageTabCol"></tr>
        </table>
      </div>

      <h4 class="mt-5 mb-1">Usage Info</h4>
      <div class="mt-3 mb-2" style="overflow: auto;">
        <table mat-table [dataSource]="usageTabDat" style="width: 100%;">
          <ng-container matColumnDef="path">
            <th mat-header-cell *matHeaderCellDef>Path</th>
            <td mat-cell *matCellDef="let u">{{ u.path }}</td>
          </ng-container>
          <ng-container matColumnDef="used">
            <th mat-header-cell *matHeaderCellDef>Used</th>
            <td mat-cell *matCellDef="let u">{{ u.usedText }}</td>
          </ng-container>
          <tr mat-header-row *matHeaderRowDef="usageTabCol"></tr>
          <tr mat-row *matRowDef="let row; columns: usageTabCol"></tr>
        </table>
      </div>
    </div>
  `,
  styles: [],
})
export class FstoreStorageComponent implements OnInit {
  storageTabCol = ["mounted", "total", "used", "available", "percent"];
  usageTabCol = ["path", "used"];

  storageTabDat: VolumnInfo[] = [];
  usageTabDat: StorageUsageInfo[] = [];

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private userService: UserService
  ) {}

  ngOnInit(): void {
    this.userService.resourceObservable.subscribe(() => {
      if (this.hasRes("fstore:server:maintenance")) {
        this.fetch();
      }
    });
  }

  hasRes(code) {
    return this.userService.hasResource(code);
  }

  fetchStorageInfo() {
    this.http.get<any>(`/fstore/storage/info`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: StorageInfo = resp.data;
        this.storageTabDat = dat.volumns;
        if (this.storageTabDat == null) {
          this.storageTabDat = [];
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

  fetchStorageUsageInfo() {
    this.http.get<any>(`/fstore/storage/usage-info`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: StorageUsageInfo[] = resp.data;
        this.usageTabDat = dat;
        if (this.usageTabDat == null) {
          this.usageTabDat = [];
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

  removeDeletedFiles() {
    this.http.post<any>(`/fstore/maintenance/remove-deleted`, null).subscribe({
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

  sanitizeStorage() {
    this.http
      .post<any>(`/fstore/maintenance/sanitize-storage`, null)
      .subscribe({
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

  fetch() {
    this.fetchStorageInfo();
    this.fetchStorageUsageInfo();
  }

  compensateMissingThumbnails() {
    this.http.post<any>(`/vfm/compensate/thumbnail`, null).subscribe({
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

  regenerateVideoThumbnails() {
    this.http
      .post<any>(`/vfm/compensate/regenerate-video-thumbnails`, null)
      .subscribe({
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
}
