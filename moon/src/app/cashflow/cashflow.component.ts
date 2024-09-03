import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpClient, HttpEvent } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { Paging, PagingController } from "src/common/paging";
import { UserService } from "../user.service";
import { isEnterKey } from "src/common/condition";
import { FormControl, FormGroup } from "@angular/forms";

export interface ListCashFlowReq {
  paging?: Paging;
  direction?: string; // Flow Direction: IN / OUT
  transTimeStart?: number; // Transaction Time Range Start
  transTimeEnd?: number; // Transaction Time Range End
  transId?: string; // Transaction ID
  category?: string; // Category Code
  minAmt?: string; // Minimum amount
}

export interface ListCashFlowRes {
  direction?: string; // Flow Direction: IN / OUT
  transTime?: number; // Transaction Time
  transId?: string; // Transaction ID
  counterparty?: string; // Counterparty of the transaction
  paymentMethod?: string; // Payment Method
  amount?: string; // Amount
  currency?: string; // Currency
  extra?: string; // Extra Information
  category?: string; // Category Code
  categoryName?: string; // Category Name
  remark?: string; // Remark
  createdAt?: number; // Create Time
}

@Component({
  selector: "app-cashflow",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Cashflows</h3>
    </div>

    <div class="m-4 shadow rounded p-3 border" *ngIf="showUploadPanel">
      <h5>Import {{ importType }} Cashflows:</h5>
      <div class="input-group input-group-lg mt-1 mb-1">
        <input
          type="file"
          class="form-control"
          (change)="onFileSelected($event.target.files)"
        />
      </div>

      <div class="d-grid gap-2 d-md-flex justify-content-md-end m-3">
        <button
          class="ml-2 mr-2"
          mat-raised-button
          [disabled]="!file"
          (click)="importCashflows()"
        >
          Upload
        </button>
      </div>
    </div>

    <div class="row row-cols-lg-auto g-3 align-items-center">
      <mat-form-field style="width: 90%;" class="mb-1 mt-3">
        <mat-label>Transaction ID</mat-label>
        <input
          matInput
          type="text"
          [(ngModel)]="fetchListReq.transId"
          (keyup)="isEnterKeyPressed($event) && fetchList()"
        />
        <button
          *ngIf="fetchListReq.transId"
          matSuffix
          aria-label="Clear"
          (click)="fetchListReq.transId = ''"
          class="btn-close"
        ></button>
      </mat-form-field>
    </div>

    <div class="row row-cols-lg-auto g-3 align-items-center">
      <div class="col">
        <mat-form-field>
          <mat-label>Category</mat-label>
          <mat-select
            (valueChange)="onCategorySelected($event)"
            [value]="fetchListReq.category"
          >
            <mat-option
              [value]="option.value"
              *ngFor="let option of [{ name: 'Wechat', value: 'WECHAT' }]"
            >
              {{ option.name }}
            </mat-option>
          </mat-select>
        </mat-form-field>
      </div>
      <!-- </div> -->

      <!-- <div class="row row-cols-lg-auto g-3 align-items-center"> -->
      <div class="col">
        <mat-form-field>
          <mat-label>Direction</mat-label>
          <mat-select
            (valueChange)="onDirectionSelected($event)"
            [value]="fetchListReq.direction"
          >
            <mat-option
              [value]="option.value"
              *ngFor="
                let option of [
                  { name: 'In', value: 'IN' },
                  { name: 'Out', value: 'OUT' }
                ]
              "
            >
              {{ option.name }}
            </mat-option>
          </mat-select>
        </mat-form-field>
      </div>

      <div class="col">
        <mat-form-field style="width: 150px" class="mb-1 mt-1">
          <mat-label>Minimum Amount</mat-label>
          <input
            matInput
            type="text"
            [(ngModel)]="fetchListReq.minAmt"
            (keyup)="isEnterKeyPressed($event) && fetchList()"
          />
          <button
            *ngIf="fetchListReq.transId"
            matSuffix
            aria-label="Clear"
            (click)="fetchListReq.minAmt = null"
            class="btn-close"
          ></button>
        </mat-form-field>
      </div>
    </div>

    <mat-form-field appearance="fill">
      <mat-label>Transaction Time</mat-label>
      <mat-date-range-input [formGroup]="range" [rangePicker]="picker">
        <input
          matStartDate
          (dateChange)="fetchList()"
          formControlName="start"
          placeholder="Start date"
        />
        <input
          matEndDate
          (dateChange)="fetchList()"
          formControlName="end"
          placeholder="End date"
        />
      </mat-date-range-input>
      <mat-datepicker-toggle matSuffix [for]="picker"></mat-datepicker-toggle>
      <mat-date-range-picker #picker></mat-date-range-picker>
    </mat-form-field>

    <div class="d-grid gap-2 d-md-flex justify-content-md-end mb-3">
      <button mat-raised-button class="m-2" (click)="showWechatImport()">
        Import Wechat Cashflows
      </button>
      <button mat-raised-button class="m-2" (click)="fetchList()">Fetch</button>
      <button mat-raised-button class="m-2" (click)="reset()">Reset</button>
    </div>

    <div class="mt-3 mb-2" style="overflow: auto;">
      <table mat-table [dataSource]="tabdat" style="width: 100%;">
        <ng-container matColumnDef="transId">
          <th mat-header-cell *matHeaderCellDef>Transaction ID</th>
          <td mat-cell *matCellDef="let u">{{ u.transId }}</td>
        </ng-container>

        <ng-container matColumnDef="transTime">
          <th mat-header-cell *matHeaderCellDef>Transaction Time</th>
          <td mat-cell *matCellDef="let u">
            {{ u.transTime | date : "yyyy-MM-dd HH:mm:ss" }}
          </td>
        </ng-container>

        <!-- <ng-container matColumnDef="direction">
          <th mat-header-cell *matHeaderCellDef>Direction</th>
          <td mat-cell *matCellDef="let u">
            <span *ngIf="u.direction == 'IN'" class="greenspan">
              {{ u.direction }}
            </span>
            <span *ngIf="u.direction != 'IN'" class="redspan">
              {{ u.direction }}
            </span>
          </td>
        </ng-container> -->

        <ng-container matColumnDef="counterparty">
          <th mat-header-cell *matHeaderCellDef>Counteryparty</th>
          <td mat-cell *matCellDef="let u">{{ u.counterparty }}</td>
        </ng-container>

        <ng-container matColumnDef="amount">
          <th mat-header-cell *matHeaderCellDef>Amount</th>
          <td
            mat-cell
            *matCellDef="let u"
            [ngClass]="u.amount.startsWith('-') ? 'redtext' : 'greentext'"
          >
            {{ u.amount }}
          </td>
        </ng-container>

        <ng-container matColumnDef="currency">
          <th mat-header-cell *matHeaderCellDef>Currency</th>
          <td mat-cell *matCellDef="let u">{{ u.currency }}</td>
        </ng-container>

        <ng-container matColumnDef="category">
          <th mat-header-cell *matHeaderCellDef>Category</th>
          <td mat-cell *matCellDef="let u">{{ u.categoryName }}</td>
        </ng-container>

        <ng-container matColumnDef="paymentMethod">
          <th mat-header-cell *matHeaderCellDef>Payment Method</th>
          <td mat-cell *matCellDef="let u">{{ u.paymentMethod }}</td>
        </ng-container>

        <ng-container matColumnDef="remark">
          <th mat-header-cell *matHeaderCellDef>Remark</th>
          <td mat-cell *matCellDef="let u">{{ u.remark }}</td>
        </ng-container>

        <ng-container matColumnDef="createdAt">
          <th mat-header-cell *matHeaderCellDef>Create Time</th>
          <td mat-cell *matCellDef="let u">
            {{ u.createdAt | date : "yyyy-MM-dd HH:mm:ss" }}
          </td>
        </ng-container>

        <ng-container matColumnDef="operation">
          <th mat-header-cell *matHeaderCellDef>Operation</th>
          <td mat-cell *matCellDef="let u">
            <button mat-raised-button (click)="popToRemove(u.transId, u.name)">
              Remove
            </button>
          </td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="tabcol"></tr>
        <tr mat-row *matRowDef="let row; columns: tabcol"></tr>
      </table>
    </div>

    <app-controlled-paginator
      (controllerReady)="onPagingControllerReady($event)"
    ></app-controlled-paginator>
  `,
  styles: [],
})
export class CashflowComponent implements OnInit {
  tabdat: ListCashFlowRes[] = [];
  tabcol = [
    // "direction",
    "transTime",
    "transId",
    "counterparty",
    "amount",
    "currency",
    "category",
    "paymentMethod",
    "remark",
    "createdAt",
  ];
  fetchListReq: ListCashFlowReq = {};
  pagingController: PagingController;
  showUploadPanel: boolean = false;
  file: File = null;
  importType: string;
  isEnterKeyPressed = isEnterKey;

  range = new FormGroup({
    start: new FormControl(),
    end: new FormControl(),
  });

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private userService: UserService
  ) {}

  ngOnInit(): void {}

  popToRemove(transId, name) {}

  fetchList() {
    this.fetchListReq.paging = this.pagingController.paging;

    this.fetchListReq.transTimeStart = this.range.value.start
      ? this.range.value.start.getTime()
      : null;

    this.fetchListReq.transTimeEnd = this.range.value.end
      ? this.range.value.end.getTime()
      : null;

    this.http
      .post<any>(`/acct/open/api/v1/cashflow/list`, this.fetchListReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          let paging: Paging = resp.data.paging;
          let payload: ListCashFlowRes[] = resp.data.payload;
          this.tabdat = payload;
          if (this.tabdat == null) {
            this.tabdat = [];
          }
          for (let r of this.tabdat) {
            if (r.direction == "OUT") {
              r.amount = "-" + r.amount;
            }
          }
          this.pagingController.onTotalChanged(paging);
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  importCashflows() {
    let type = this.importType;
    if (!type) {
      return;
    }
    if (!this.file) {
      this.snackBar.open("Please select file", "ok", { duration: 3000 });
      return;
    }

    if (type == "Wechat") {
      this.http
        .post<any>(`/acct/open/api/v1/cashflow/import/wechat`, this.file)
        .subscribe({
          complete: () => {
            this.file = null;
            this.showUploadPanel = false;
            this.importType = "";
            this.snackBar.open("Uploaded, this may take a while", "ok", {
              duration: 3000,
            });
            setTimeout(() => this.fetchList(), 1000);
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }

  showWechatImport() {
    this.showUploadPanel = !this.showUploadPanel;
    if (this.showUploadPanel) {
      this.importType = "Wechat";
    } else {
      this.importType = "";
    }
  }

  onFileSelected(files: File[]) {
    if (files == null || files.length < 1) {
      this.snackBar.open("Please select file", "ok", { duration: 3000 });
      return;
    }
    this.file = files[0];
  }

  reset() {
    this.fetchListReq = {};
    this.range.reset();
    this.pagingController.firstPage();
    this.fetchList();
  }

  onDirectionSelected(dir) {
    this.fetchListReq.direction = dir;
    this.fetchList();
  }

  onCategorySelected(cat) {
    this.fetchListReq.category = cat;
    this.fetchList();
  }
}
