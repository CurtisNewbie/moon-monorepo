import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatSelectModule } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";

import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { I18n } from "../i18n.service";
import { NavigationService } from "../navigation.service";
import { TrlPipe } from "../trl.pipe";
import { DroneTaskComponent } from "./drone-task.component";

const LIST_URL = "/drone/open/api/list-task";

describe("DroneTaskComponent search", () => {
  let fixture: ComponentFixture<DroneTaskComponent>;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [
        DroneTaskComponent,
        ControlledPaginatorComponent,
        TrlPipe,
      ],
      imports: [
        HttpClientTestingModule,
        NoopAnimationsModule,
        FormsModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatSelectModule,
        MatInputModule,
      ],
      providers: [
        { provide: MatSnackBar, useValue: { open: () => undefined } },
        { provide: MatDialog, useValue: {} },
        { provide: I18n, useValue: { trl: (_m: string, k: string) => k } },
        { provide: NavigationService, useValue: {} },
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { queryParams: {} } },
        },
        { provide: Env, useValue: { isMobile: () => false } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    });
    httpMock = TestBed.inject(HttpTestingController);
    fixture = TestBed.createComponent(DroneTaskComponent);
    fixture.detectChanges();

    // the paginator fires the first page fetch on init
    flushAll();
  });

  afterEach(() => {
    httpMock.verify();
  });

  function flushAll(): void {
    httpMock
      .match(LIST_URL)
      .forEach((r) => r.flush({ data: { payload: [], paging: { total: 0 } } }));
    httpMock.match("/drone/open/api/list-platforms").forEach((r) => r.flush({ data: [] }));
    httpMock
      .match("/drone/open/api/list-last-page-platforms")
      .forEach((r) => r.flush({ data: [] }));
  }

  /** The search box whose label reads as the given i18n key. */
  function searchBox(labelKey: string): HTMLInputElement {
    const fields = Array.from(
      fixture.nativeElement.querySelectorAll("mat-form-field")
    ) as HTMLElement[];
    const field = fields.find(
      (f) => (f.querySelector("mat-label")?.textContent || "").trim() === labelKey
    );
    if (!field) {
      throw new Error(`no search box labelled '${labelKey}'`);
    }
    return field.querySelector("input") as HTMLInputElement;
  }

  function fillAndSearch(labelKey: string, value: string): void {
    const input = searchBox(labelKey);
    input.value = value;
    input.dispatchEvent(new Event("input"));
    input.dispatchEvent(new KeyboardEvent("keyup", { key: "Enter" }));
    fixture.detectChanges();
  }

  function lastListBody(): any {
    const reqs = httpMock.match(LIST_URL);
    expect(reqs.length).toBe(1);
    reqs[0].flush({ data: { payload: [], paging: { total: 0 } } });
    return reqs[0].request.body;
  }

  it("searches by directory name when that box is filled and Enter is pressed", () => {
    fillAndSearch("searchDirName", "my archive");

    expect(lastListBody().dirName).toBe("my archive");
  });

  it("sends no directory name filter when that box is empty", () => {
    fillAndSearch("searchUrl", "https://example.com/a");

    const body = lastListBody();
    expect(body.url).toBe("https://example.com/a");
    expect(body.dirName).toBeUndefined();
  });

  it("sends both filters when both boxes are filled", () => {
    fillAndSearch("searchUrl", "https://example.com/a");
    lastListBody();
    fillAndSearch("searchDirName", "my archive");

    const body = lastListBody();
    expect(body.url).toBe("https://example.com/a");
    expect(body.dirName).toBe("my archive");
  });
});
