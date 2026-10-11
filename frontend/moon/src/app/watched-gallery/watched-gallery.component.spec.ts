import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
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
import { WatchedGalleryComponent } from "./watched-gallery.component";

const LIST_URL = "/drone/open/api/list-watched-gallery";

describe("WatchedGalleryComponent paginator fetch", () => {
  let fixture: ComponentFixture<WatchedGalleryComponent>;
  let httpMock: HttpTestingController;

  function configure(isMobile: boolean): void {
    TestBed.configureTestingModule({
      declarations: [
        WatchedGalleryComponent,
        ControlledPaginatorComponent,
        TrlPipe,
      ],
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
        { provide: MatSnackBar, useValue: { open: () => undefined } },
        { provide: I18n, useValue: { trl: (_m: string, k: string) => k } },
        { provide: NavigationService, useValue: {} },
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { queryParams: { taskId: "" } } },
        },
        { provide: Env, useValue: { isMobile: () => isMobile } },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    });
    httpMock = TestBed.inject(HttpTestingController);
  }

  afterEach(() => {
    httpMock.verify();
  });

  it("issues exactly one list request on mobile, using the paginator's page limit (5)", () => {
    configure(true);
    fixture = TestBed.createComponent(WatchedGalleryComponent);
    fixture.detectChanges();

    const reqs = httpMock.match(LIST_URL);
    reqs.forEach((r) =>
      r.flush({ data: { payload: [], paging: { total: 0 } } })
    );

    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body.paging.limit).toBe(5);
  });

  it("issues exactly one list request on desktop, using the paginator's page limit (10)", () => {
    configure(false);
    fixture = TestBed.createComponent(WatchedGalleryComponent);
    fixture.detectChanges();

    const reqs = httpMock.match(LIST_URL);
    reqs.forEach((r) =>
      r.flush({ data: { payload: [], paging: { total: 0 } } })
    );

    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body.paging.limit).toBe(10);
  });
});
