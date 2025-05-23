import { MatSnackBar } from "@angular/material/snack-bar";
import { HttpClient } from "@angular/common/http";
import { Component, OnInit, ViewChild } from "@angular/core";
import { Paging} from "src/common/paging";
import { isEnterKey } from "src/common/condition";
import { FormControl, FormGroup } from "@angular/forms";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

export interface ApiPlotStatisticsReq {
  startTime?: number; // Start time
  endTime?: number; // End time
  aggType?: string; // Aggregation Type.
  currency?: string; // Currency
}

export interface ApiPlotStatisticsRes {
  aggRange?: string; // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  aggValue?: string; // Aggregation Value.
}

export interface ApiListStatisticsReq {
  paging?: Paging;
  aggType?: string; // Aggregation Type.
  aggRange?: string; // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  currency?: string; // Currency
}

export interface ApiListStatisticsRes {
  aggType?: string; // Aggregation Type.
  aggRange?: string; // Aggregation Range. The corresponding year (YYYY), month (YYYYMM), sunday of the week (YYYYMMDD).
  aggValue?: string; // Aggregation Value.
  currency?: string; // Currency
}

@Component({
  selector: "app-cashflow-statistics",
  template: `
    <div>
      <h3 class="mt-2 mb-3">Cashflow Statistics</h3>
    </div>

    <plotly-plot
      class="darkmode"
      [data]="graph.data"
      [layout]="graph.layout"
      [useResizeHandler]="true"
      [style]="{ position: 'relative', width: '100%', height: '100%' }"
    ></plotly-plot>

    <div class="mt-4 d-flex flex-wrap justify-content-between">
      <div class="d-flex flex-wrap align-items-center gap-3">
        <mat-form-field>
          <mat-label>Plot Date Range</mat-label>
          <mat-date-range-input [formGroup]="range" [rangePicker]="picker">
            <input
              matStartDate
              (dateChange)="fetchPlots()"
              formControlName="start"
              placeholder="Start date"
            />
            <input
              matEndDate
              (dateChange)="fetchPlots()"
              formControlName="end"
              placeholder="End date"
            />
          </mat-date-range-input>
          <mat-datepicker-toggle
            matSuffix
            [for]="picker"
          ></mat-datepicker-toggle>
          <mat-date-range-picker #picker></mat-date-range-picker>
        </mat-form-field>
        <mat-form-field>
          <mat-label>Type</mat-label>
          <mat-select
            (valueChange)="onAggTypeSelected($event)"
            [value]="listReq.aggType"
          >
            <mat-option
              [value]="option.value"
              *ngFor="
                let option of [
                  { name: 'Yearly', value: 'YEARLY' },
                  { name: 'Monthly', value: 'MONTHLY' },
                  { name: 'Weekly', value: 'WEEKLY' }
                ]
              "
            >
              {{ option.name }}
            </mat-option>
          </mat-select>
        </mat-form-field>
        <mat-form-field>
          <mat-label>Currency</mat-label>
          <mat-select
            (valueChange)="onCurrencySelected($event)"
            [value]="listReq.currency"
          >
            <mat-option [value]="option" *ngFor="let option of currencies">
              {{ option }}
            </mat-option>
          </mat-select>
        </mat-form-field>
      </div>

      <div>
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

    <div class="mt-3 mb-2" style="overflow: auto;">
      <table mat-table [dataSource]="dat" style="width: 100%;">
        <ng-container matColumnDef="aggType">
          <th mat-header-cell *matHeaderCellDef>Type</th>
          <td mat-cell *matCellDef="let u">{{ u.aggType }}</td>
        </ng-container>

        <ng-container matColumnDef="aggRange">
          <th mat-header-cell *matHeaderCellDef>Range</th>
          <td mat-cell *matCellDef="let u">{{ u.aggRange }}</td>
        </ng-container>

        <ng-container matColumnDef="aggValue">
          <th mat-header-cell *matHeaderCellDef>Amount</th>
          <td
            mat-cell
            *matCellDef="let u"
            [ngClass]="u.aggValue.startsWith('-') ? 'redtext' : 'greentext'"
          >
            {{ u.aggValue }}
          </td>
        </ng-container>

        <ng-container matColumnDef="currency">
          <th mat-header-cell *matHeaderCellDef>Currency</th>
          <td mat-cell *matCellDef="let u">{{ u.currency }}</td>
        </ng-container>

        <tr mat-header-row *matHeaderRowDef="tabcol"></tr>
        <tr mat-row *matRowDef="let row; columns: tabcol"></tr>
      </table>
    </div>

    <app-controlled-paginator
      (pageChanged)="fetchList()"
    ></app-controlled-paginator>
  `,
  styles: [],
})
export class CashflowStatisticsComponent implements OnInit {
  tabcol = ["aggType", "aggRange", "aggValue", "currency"];
  isEnterKey = isEnterKey;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  dat: ApiListStatisticsRes[] = [];
  plotDat: ApiPlotStatisticsRes[] = [];
  currencies = [];
  listReq: ApiListStatisticsReq = {
    aggType: "WEEKLY",
  };
  plotReq: ApiPlotStatisticsReq = {};

  range = new FormGroup({
    start: new FormControl(
      new Date(new Date().setFullYear(new Date().getFullYear() - 1))
    ),
    end: new FormControl(new Date()),
  });

  public graph = {
    data: [
      {
        x: [],
        y: [],
        type: "scatter",
        mode: "lines",
      },
    ],
    layout: {
      height: 350,
      xaxis: {
        labelalias: {},
        title: "Date Range",
      },
      yaxis: {
        title: "Cashflow",
      },
      // title: "Cashflow Statistics",
    },
  };

  constructor(private snackBar: MatSnackBar, private http: HttpClient) {}

  ngOnInit(): void {
    this.pagingController.PAGE_LIMIT_OPTIONS = [5, 10, 30, 50];
    this.pagingController.setPageLimit(5);
  }

  fetchList() {
    this.listReq.paging = this.pagingController.paging;
    this.http
      .post<any>(`/acct/open/api/v1/cashflow/list-statistics`, this.listReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.dat = resp.data.payload;
          this.pagingController.onTotalChanged(resp.data.paging);
          if (this.dat == null) {
            this.dat = [];
          }
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });

    this.fetchPlots();
    this.fetchCurrencies();
  }

  reset() {
    this.listReq = {
      aggType: "WEEKLY",
    };
    if (!this.pagingController.firstPage()) {
      this.fetchList();
    }
  }

  onCurrencySelected(currency) {
    this.listReq.currency = currency;
    if (!this.pagingController.firstPage()) {
      this.fetchList();
    }
  }

  onAggTypeSelected(aggType) {
    this.listReq.aggType = aggType;
    if (!this.pagingController.firstPage()) {
      this.fetchList();
    }
    this.fetchPlots();
  }

  fetchCurrencies() {
    this.http.get<any>(`/acct/open/api/v1/cashflow/list-currency`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: string[] = resp.data;
        if (dat == null) {
          dat = [];
        }
        this.currencies = dat;
        if (this.currencies.length == 1 && !this.listReq.currency) {
          this.listReq.currency = this.currencies[0];
          this.fetchPlots();
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

  fetchPlots() {
    if (
      !this.listReq.currency ||
      !this.range.value.start ||
      !this.range.value.end
    ) {
      return;
    }

    this.plotReq.currency = this.listReq.currency;
    this.plotReq.aggType = this.listReq.aggType;

    this.range.value.start.setHours(0, 0, 0, 0);
    this.range.value.end.setHours(23, 59, 59, 99);
    this.plotReq.startTime = this.range.value.start.getTime();
    this.plotReq.endTime = this.range.value.end.getTime();

    this.http
      .post<any>(`/acct/open/api/v1/cashflow/plot-statistics`, this.plotReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          let dat: ApiPlotStatisticsRes[] = resp.data;
          this.plotDat = dat;
          if (this.plotDat == null) {
            this.plotDat = [];
          }
          let x = [];
          let y = [];
          this.graph.data[0].x = x;
          this.graph.data[0].y = y;
          let i = 0;
          for (let v of this.plotDat) {
            x.push(i);
            y.push(v.aggValue);
            this.graph.layout.xaxis.labelalias[i] = v.aggRange;
            i += 1;
          }
          console.log(x);
          console.log(y);
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
