import { Component, OnInit } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpClient } from "@angular/common/http";
import { UserService } from "../user.service";

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
      <h3 class="mt-2 mb-3">Mini-Fstore Storage Info</h3>
    </div>

    <div *ngIf="hasRes('fstore:maintenance')" class="d-flex flex-wrap justify-content-end">
      <button
        mat-button
        (click)="sanitizeStorage()"
        matTooltip="Delete dangling files that are not uploaded to server"
      >
        Sanitize Storage
      </button>
      <button
        mat-button
        (click)="removeDeletedFiles()"
        matTooltip="Prune files that are marked deleted. During server maintenance, file uploading is rejected"
      >
        Prune Deleted
      </button>
    </div>

    <div class="mt-3 mb-2" style="overflow: auto;">
      <table mat-table [dataSource]="tabdat" style="width: 100%;">
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

        <tr mat-header-row *matHeaderRowDef="tabcol"></tr>
        <tr mat-row *matRowDef="let row; columns: tabcol"></tr>
      </table>
    </div>
  `,
  styles: [],
})
export class FstoreStorageComponent implements OnInit {
  tabcol = ["mounted", "total", "used", "available", "percent"];
  tabdat: VolumnInfo[] = [];
  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private userService: UserService
  ) {}

  ngOnInit(): void {
    this.fetchStorageInfo();
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
        this.tabdat = dat.volumns;
        if (this.tabdat == null) {
          this.tabdat = [];
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
}
