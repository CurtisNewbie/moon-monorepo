import { Component, OnDestroy, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpClient } from "@angular/common/http";
import { UserService } from "../user.service";
import { timer } from "rxjs";

export interface MaintenanceStatus {
  underMaintenance?: boolean;
}

export interface StorageUsageInfo {
  type?: string;
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
          [disabled]="vfmMaintenance"
          (click)="compensateMissingThumbnails()"
          matTooltip="Generate missing thumbnails"
        >
          Compensate Missing Thumbnails
        </button>
        <button
          mat-raised-button
          [disabled]="vfmMaintenance"
          (click)="regenerateVideoThumbnails()"
          matTooltip="Regenerate video thumbnails"
        >
          Regenerate Video Thumbnails
        </button>
      </ng-container>

      <ng-container *ngIf="hasRes('fstore:server:maintenance')">
        <button
          mat-raised-button
          [disabled]="fstoreMaintenance"
          (click)="sanitizeStorage()"
          matTooltip="Delete dangling files that are not uploaded to server"
        >
          Sanitize Storage
        </button>
        <button
          mat-raised-button
          [disabled]="fstoreMaintenance"
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
          <ng-container matColumnDef="type">
            <th mat-header-cell *matHeaderCellDef>Type</th>
            <td mat-cell *matCellDef="let u">{{ u.type }}</td>
          </ng-container>
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
export class FstoreStorageComponent implements OnInit, OnDestroy {
  storageTabCol = ["mounted", "total", "used", "available", "percent"];
  usageTabCol = ["type", "path", "used"];

  storageTabDat: VolumnInfo[] = [];
  usageTabDat: StorageUsageInfo[] = [];
  fstoreMaintenance: boolean = false;
  vfmMaintenance: boolean = false;

  checkMaintenanceTimerSub = timer(0, 3000).subscribe(() => {
    this.checkFstoreMaintenanceStatus();
    this.checkVfmMaintenanceStatus();
  });

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private userService: UserService
  ) {}

  ngOnDestroy(): void {
    this.checkMaintenanceTimerSub.unsubscribe();
  }

  ngOnInit(): void {
    this.userService.resourceObservable.subscribe(() => {
      if (this.hasRes("fstore:server:maintenance")) {
        this.fetch();
      }
    });
    this.checkFstoreMaintenanceStatus();
    this.checkVfmMaintenanceStatus();
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
        this.fstoreMaintenance = true;
        this.snackBar.open("Request success, make take a while", "ok", {
          duration: 1500,
        });
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
          this.fstoreMaintenance = true;
          this.snackBar.open("Request success, make take a while", "ok", {
            duration: 1500,
          });
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
    this.checkFstoreMaintenanceStatus();
    this.checkVfmMaintenanceStatus();
  }

  compensateMissingThumbnails() {
    this.http.post<any>(`/vfm/compensate/thumbnail`, null).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.vfmMaintenance = true;
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
          this.vfmMaintenance = true;
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  checkFstoreMaintenanceStatus() {
    this.http.get<any>(`/fstore/maintenance/status`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: MaintenanceStatus = resp.data;
        this.fstoreMaintenance = dat.underMaintenance;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }

  checkVfmMaintenanceStatus() {
    this.http.get<any>(`/vfm/maintenance/status`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: MaintenanceStatus = resp.data;
        this.vfmMaintenance = dat.underMaintenance;
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
