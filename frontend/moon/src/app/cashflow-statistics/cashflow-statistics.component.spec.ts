import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatNativeDateModule } from "@angular/material/core";
import { MatDatepickerModule } from "@angular/material/datepicker";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatSelectModule } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { CashflowStatisticsComponent } from "./cashflow-statistics.component";

const LIST_URL = "/acct/open/api/v1/cashflow/list-statistics";
const CURRENCY_URL = "/acct/open/api/v1/cashflow/list-currency";

describe("CashflowStatisticsComponent paginator default limit", () => {
  let fixture: ComponentFixture<CashflowStatisticsComponent>;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [CashflowStatisticsComponent, ControlledPaginatorComponent],
      imports: [
        HttpClientTestingModule,
        NoopAnimationsModule,
        FormsModule,
        ReactiveFormsModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatSelectModule,
        MatInputModule,
      ],
      providers: [
        { provide: MatSnackBar, useValue: { open: () => undefined } },
        // Force the non-mobile path so the paginator keeps its class default
        // (10) unless ngOnInit actually applies the configured options.
        { provide: Env, useValue: { isMobile: () => false } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    });
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it("applies the default page limit (5) and does not throw during init", () => {
    fixture = TestBed.createComponent(CashflowStatisticsComponent);

    expect(() => fixture.detectChanges()).not.toThrow();

    expect(fixture.componentInstance.pagingController.paging.limit).toBe(5);

    httpMock
      .match(LIST_URL)
      .forEach((r) => r.flush({ data: { payload: [], paging: { total: 0 } } }));
    httpMock.match(CURRENCY_URL).forEach((r) => r.flush({ data: [] }));
  });

  it("sends the first list request with paging.limit 5", () => {
    fixture = TestBed.createComponent(CashflowStatisticsComponent);
    fixture.detectChanges();

    const reqs = httpMock.match(LIST_URL);
    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body.paging.limit).toBe(5);

    reqs.forEach((r) => r.flush({ data: { payload: [], paging: { total: 0 } } }));
    httpMock.match(CURRENCY_URL).forEach((r) => r.flush({ data: [] }));
  });
});
