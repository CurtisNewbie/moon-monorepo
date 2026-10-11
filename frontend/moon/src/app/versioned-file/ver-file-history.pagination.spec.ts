import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogRef,
} from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatSelectModule } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { FileInfoService } from "../file-info.service";
import { NavigationService } from "../navigation.service";
import { VerFileHistoryComponent } from "./versioned-file.component";

const LIST_URL = "vfm/open/api/versioned-file/history";
const ACCUMULATED_SIZE_URL = "vfm/open/api/versioned-file/accumulated-size";

describe("VerFileHistoryComponent paginator default limit", () => {
  let fixture: ComponentFixture<VerFileHistoryComponent>;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [VerFileHistoryComponent, ControlledPaginatorComponent],
      imports: [
        HttpClientTestingModule,
        NoopAnimationsModule,
        FormsModule,
        ReactiveFormsModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
      ],
      providers: [
        { provide: MatDialogRef, useValue: { close: () => undefined } },
        { provide: MAT_DIALOG_DATA, useValue: { verFileId: "test-file" } },
        { provide: FileInfoService, useValue: {} },
        { provide: MatDialog, useValue: {} },
        { provide: NavigationService, useValue: {} },
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
    fixture = TestBed.createComponent(VerFileHistoryComponent);

    expect(() => fixture.detectChanges()).not.toThrow();

    expect(fixture.componentInstance.pagingController.paging.limit).toBe(5);

    httpMock
      .match(LIST_URL)
      .forEach((r) => r.flush({ data: { payload: [], paging: { total: 0 } } }));
    httpMock
      .match(ACCUMULATED_SIZE_URL)
      .forEach((r) => r.flush({ data: { sizeInBytes: 0 } }));
  });

  it("sends the first history list request with paging.limit 5", () => {
    fixture = TestBed.createComponent(VerFileHistoryComponent);
    fixture.detectChanges();

    const reqs = httpMock.match(LIST_URL);
    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body.paging.limit).toBe(5);

    reqs.forEach((r) => r.flush({ data: { payload: [], paging: { total: 0 } } }));
    httpMock
      .match(ACCUMULATED_SIZE_URL)
      .forEach((r) => r.flush({ data: { sizeInBytes: 0 } }));
  });
});
