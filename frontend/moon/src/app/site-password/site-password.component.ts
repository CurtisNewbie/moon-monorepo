import { HttpClient } from "@angular/common/http";
import { Component, Inject, OnInit, ViewChild } from "@angular/core";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { MatSnackBar } from "@angular/material/snack-bar";
import { copyToClipboard } from "src/common/clipboard";
import { isEnterKey } from "src/common/condition";
import { ConfirmDialog } from "src/common/dialog";
import { Env } from "src/common/env-util";
import { Paging } from "src/common/paging";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface EditSitePasswordReq {
  recordId?: string;
  site?: string;
  alias?: string;
  username?: string;
  sitePassword?: string;
  loginPassword?: string;
}
export interface DecryptSitePasswordDialogData {
  title?: string;
  sitePassword?: ListSitePasswordRes;
}
export interface AddSitePasswordReq {
  site?: string;
  alias?: string;
  username?: string;
  sitePassword?: string;
  loginPassword?: string;
}
export interface ListSitePasswordReq {
  alias?: string;
  site?: string;
  username?: string;
  paging?: Paging;
}
export interface ListSitePasswordRes {
  recordId?: string;
  site?: string;
  alias?: string;
  username?: string;
  createTime?: number;
}
export interface DecryptSitePasswordReq {
  loginPassword?: string;
  recordId?: string;
}
export interface RemoveSitePasswordRes {
  recordId?: string;
}

@Component({
  selector: "edit-site-password-dialog",
  template: `
    <h1 mat-dialog-title>Edit Site Password</h1>
    <div>
      <ng-container>
        <p>Record Id: {{ data.recordId }}</p>
      </ng-container>

      <ng-container>
        <mat-form-field style="width: 100%;">
          <mat-label>Alias</mat-label>
          <input matInput [(ngModel)]="editReq.alias" />
        </mat-form-field>
        <mat-form-field style="width: 100%;">
          <mat-label>Site</mat-label>
          <input matInput [(ngModel)]="editReq.site" />
        </mat-form-field>
        <mat-form-field style="width: 100%;">
          <mat-label>New Username (optional)</mat-label>
          <input matInput [(ngModel)]="editReq.username" />
        </mat-form-field>
        <mat-form-field style="width: 100%;">
          <mat-label>New Site Password</mat-label>
          <input
            autocomplete="one-time-code"
            type="password"
            matInput
            [(ngModel)]="editReq.sitePassword"
          />
        </mat-form-field>
        <ng-container *ngIf="editReq.sitePassword">
          <mat-form-field style="width: 100%;">
            <mat-label>Login Password</mat-label>
            <input
              autocomplete="one-time-code"
              type="password"
              matInput
              [(ngModel)]="editReq.loginPassword"
            />
          </mat-form-field>
        </ng-container>
        <div class="justify-content-end d-flex">
          <button
            mat-raised-button
            class="mt-2"
            (click)="editSitePassword()"
            [mat-dialog-close]="true"
          >
            Submit
          </button>
        </div>
      </ng-container>
    </div>

    <div mat-dialog-actions>
      <button mat-button [mat-dialog-close]="true">Close</button>
    </div>
  `,
})
export class EditSitePasswordDialogComponent {
  editReq: EditSitePasswordReq;
  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    public dialogRef: MatDialogRef<
      EditSitePasswordDialogComponent,
      ListSitePasswordRes
    >,
    @Inject(MAT_DIALOG_DATA) public data: ListSitePasswordRes
  ) {
    this.editReq = {
      recordId: this.data.recordId,
      site: this.data.site,
      alias: this.data.alias,
    };
  }

  editSitePassword() {
    let req: EditSitePasswordReq = this.editReq;
    this.http
      .post<any>(`/user-vault/open/api/password/edit-site-password`, req)
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

@Component({
  selector: "site-password-decrypted-dialog",
  template: `
    <h1 mat-dialog-title>View Site Password</h1>
    <div>
      <p>{{ data.title }}</p>
      <ng-container *ngIf="decrypted">
        <mat-form-field style="width: 100%;">
          <mat-label>Decrypted:</mat-label>
          <input
            matInput
            [ngModel]="decrypted"
            readonly="true"
            (click)="copyDecrypted()"
          />
        </mat-form-field>
      </ng-container>

      <ng-container *ngIf="!decrypted">
        <mat-form-field style="width: 100%;">
          <mat-label>Login Password</mat-label>
          <input
            matInput
            autocomplete="one-time-code"
            type="password"
            [(ngModel)]="loginPasssword"
            (keyup)="
              isEnterKey($event) &&
                decryptSitePassword(data.sitePassword, loginPasssword)
            "
          />
        </mat-form-field>
        <div class="justify-content-end d-flex">
          <button
            mat-raised-button
            class="mt-2"
            (click)="decryptSitePassword(data.sitePassword, loginPasssword)"
          >
            Decrypt Password
          </button>
        </div>
      </ng-container>
    </div>

    <div mat-dialog-actions>
      <button mat-button [mat-dialog-close]="true">Close</button>
    </div>
  `,
})
export class SitePasswordDecryptedDialogComponent implements OnInit {
  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    public dialogRef: MatDialogRef<
      SitePasswordDecryptedDialogComponent,
      DecryptSitePasswordDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: DecryptSitePasswordDialogData
  ) {}

  copyDecrypted = () => {
    copyToClipboard(this.decrypted);
    this.snackBar.open("Copied to clipboard", "ok", { duration: 1000 });
  };

  loginPasssword: string = "";
  decrypted: string = "";
  isEnterKey = isEnterKey;

  decryptSitePassword(u: ListSitePasswordRes, loginPassword: string) {
    let req: DecryptSitePasswordReq = {
      loginPassword: loginPassword,
      recordId: u.recordId,
    };
    this.http
      .post<any>(`/user-vault/open/api/password/decrypt-site-password`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.decrypted = resp.data.decrypted;
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
        complete: () => {
          this.loginPasssword = "";
        },
      });
  }
  ngOnInit(): void {}
}

@Component({
  selector: "app-site-password",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Website Passwords</h3>
      <div class="justify-content-end d-flex">
        <button mat-raised-button (click)="togglePanel()">
          Add Site Passowrd
        </button>
      </div>
    </div>

    <div
      class="container bootstrap p-3 mt-3 mb-5 shadow"
      *ngIf="panelDisplayed"
    >
      <h4 class="mt-2 mb-2">Add Site Password</h4>

      <mat-form-field g style="width: 100%;">
        <mat-label>Site</mat-label>
        <input matInput type="text" [(ngModel)]="addSitePasswordReq.site" />
      </mat-form-field>

      <mat-form-field g style="width: 100%;">
        <mat-label>Alias</mat-label>
        <input matInput type="text" [(ngModel)]="addSitePasswordReq.alias" />
      </mat-form-field>

      <mat-form-field g style="width: 100%;">
        <mat-label>Site Account Name</mat-label>
        <input matInput [(ngModel)]="addSitePasswordReq.username" />
      </mat-form-field>

      <mat-form-field g style="width: 100%;">
        <mat-label>Site Password</mat-label>
        <input
          matInput
          autocomplete="one-time-code"
          type="password"
          [(ngModel)]="addSitePasswordReq.sitePassword"
        />
      </mat-form-field>

      <mat-form-field g style="width: 100%;">
        <mat-label>Login Password</mat-label>
        <input
          autocomplete="one-time-code"
          matInput
          type="password"
          [(ngModel)]="addSitePasswordReq.loginPassword"
        />
      </mat-form-field>

      <div class="justify-content-end d-flex">
        <button mat-raised-button class="mt-2" (click)="addSitePassword()">
          Submit
        </button>
      </div>
    </div>

    <div>
      <mat-form-field g style="width: 100%;" class="">
        <mat-label>Alias:</mat-label>
        <input
          matInput
          type="text"
          [(ngModel)]="listSitePasswordReq.alias"
          (keyup)="isEnter($event) && fetchList()"
        />
        <button
          *ngIf="listSitePasswordReq.alias"
          matSuffix
          aria-label="Clear"
          (click)="listSitePasswordReq.alias = ''"
          class="btn-close"
        ></button>
      </mat-form-field>

      <mat-form-field g style="width: 100%;" class="">
        <mat-label>Site:</mat-label>
        <input
          matInput
          type="text"
          [(ngModel)]="listSitePasswordReq.site"
          (keyup)="isEnter($event) && fetchList()"
        />
        <button
          *ngIf="listSitePasswordReq.site"
          matSuffix
          aria-label="Clear"
          (click)="listSitePasswordReq.site = ''"
          class="btn-close"
        ></button>
      </mat-form-field>

      <mat-form-field g style="width: 100%;" class="">
        <mat-label>Site Account Name:</mat-label>
        <input
          matInput
          type="text"
          autocomplete="one-time-code"
          [(ngModel)]="listSitePasswordReq.username"
          (keyup)="isEnter($event) && fetchList()"
        />
        <button
          *ngIf="listSitePasswordReq.username"
          matSuffix
          aria-label="Clear"
          (click)="listSitePasswordReq.username = ''"
          class="btn-close"
        ></button>
      </mat-form-field>

      <div class="d-grid gap-2 d-flex justify-content-end mb-3">
        <button
          mat-icon-button
          class="m-1 icon-button-large"
          (click)="fetchList()"
        >
          <i class="bi bi-arrow-clockwise"></i>
        </button>
        <button mat-icon-button class="m-1 icon-button-large" (click)="reset()">
          <i class="bi bi-slash-circle"></i>
        </button>
      </div>
    </div>

    <div class="mt-3 mb-5">
      <table
        mat-table
        [dataSource]="tab"
        class="mat-elevation-z8 mb-4"
        style="width: 100%;"
        multiTemplateDataRows
      >
        <ng-container matColumnDef="recordId">
          <th mat-header-cell *matHeaderCellDef>Record ID</th>
          <td mat-cell *matCellDef="let u">{{ u.recordId }}</td>
        </ng-container>
        <ng-container matColumnDef="alias">
          <th mat-header-cell *matHeaderCellDef>Alias</th>
          <td mat-cell *matCellDef="let u">{{ u.alias }}</td>
        </ng-container>
        <ng-container matColumnDef="site">
          <th mat-header-cell *matHeaderCellDef>Site</th>
          <td mat-cell *matCellDef="let u">{{ u.site }}</td>
        </ng-container>
        <ng-container matColumnDef="username">
          <th mat-header-cell *matHeaderCellDef>Username</th>
          <td mat-cell *matCellDef="let u">{{ u.username }}</td>
        </ng-container>
        <ng-container matColumnDef="createTime">
          <th mat-header-cell *matHeaderCellDef>Create Time</th>
          <td mat-cell *matCellDef="let u">
            {{ u.createTime | date : "yyyy-MM-dd HH:mm:ss" }}
          </td>
        </ng-container>
        <ng-container matColumnDef="operation">
          <th mat-header-cell *matHeaderCellDef><b>Operation</b></th>
          <td mat-cell *matCellDef="let u">
            <button
              class="small-btn m-2"
              mat-raised-button
              (click)="$event.stopPropagation() || removeSitePassword(u)"
            >
              Remove
            </button>

            <button
              class="small-btn m-2"
              mat-raised-button
              (click)="$event.stopPropagation() || edit(u)"
            >
              Edit
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="col"></tr>
        <tr
          mat-row
          *matRowDef="let row; columns: col"
          (click)="$event.stopPropagation() || preview(row)"
        ></tr>
      </table>

      <app-controlled-paginator
        (pageChanged)="fetchList()"
      ></app-controlled-paginator>
    </div>
  `,
  styles: [],
})
export class SitePasswordComponent implements OnInit {
  readonly col = this.env.isMobile()
    ? ["alias", "site", "username"]
    : ["recordId", "alias", "site", "username", "createTime", "operation"];

  constructor(
    public env: Env,
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private dialog: MatDialog,
    private confirmDialog: ConfirmDialog
  ) {}

  panelDisplayed: boolean = false;
  addSitePasswordReq: AddSitePasswordReq = {};
  listSitePasswordReq: ListSitePasswordReq = {};
  tab: ListSitePasswordRes[] = [];
  isEnter = isEnterKey;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  ngOnInit(): void {}

  togglePanel() {
    this.panelDisplayed = !this.panelDisplayed;
    this.addSitePasswordReq = {};
  }

  addSitePassword() {
    this.http
      .post<any>(
        `/user-vault/open/api/password/add-site-password`,
        this.addSitePasswordReq
      )
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.addSitePasswordReq = {};
          this.panelDisplayed = false;
          this.fetchList();
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  fetchList() {
    this.listSitePasswordReq.paging = this.pagingController.paging;
    this.http
      .post<any>(
        `/user-vault/open/api/password/list-site-passwords`,
        this.listSitePasswordReq
      )
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.pagingController.onTotalChanged(resp.data.paging);
          this.tab = resp.data.payload;
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  reset() {
    this.listSitePasswordReq = {};
    if (!this.pagingController.atFirstPage()) {
      this.pagingController.firstPage();
    } else {
      this.fetchList();
    }
  }

  preview(u: ListSitePasswordRes) {
    let n = u.site;
    if (!n) {
      n = u.alias;
    }
    this.dialog.open(SitePasswordDecryptedDialogComponent, {
      data: {
        title: `View Password for '${n}'`,
        sitePassword: u,
      },
      width: "400px",
    });
  }

  edit(u: ListSitePasswordRes) {
    this.dialog
      .open(EditSitePasswordDialogComponent, {
        data: { ...u },
        width: "600px",
      })
      .afterClosed()
      .subscribe(() => {
        this.fetchList();
      });
  }

  removeSitePassword(u: ListSitePasswordRes) {
    let n = u.site;
    if (!n) {
      n = u.alias;
    }
    this.confirmDialog.show(
      `Remove password for ${n}?`,
      [
        `Are you sure you want to remove password for ${n}?`,
        "Result cannot be reverted once you remove it.",
      ],
      () => {
        let req: RemoveSitePasswordRes = { recordId: u.recordId };
        this.http
          .post<any>(`/user-vault/open/api/password/remove-site-password`, req)
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.fetchList();
            },
            error: (err) => {
              console.log(err);
              this.snackBar.open("Request failed, unknown error", "ok", {
                duration: 3000,
              });
            },
          });
      }
    );
  }
}
