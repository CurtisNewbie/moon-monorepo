import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatPaginatorModule } from "@angular/material/paginator";
import { MatSelectModule } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";

import { ConfirmDialog } from "src/common/dialog";
import { Env } from "src/common/env-util";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { I18n } from "../i18n.service";
import { NavigationService } from "../navigation.service";
import { TrlPipe } from "../trl.pipe";
import { FapeGalleryComponent } from "./fape-gallery.component";

const LIST_URL = "/drone/open/api/fape/gallery/list";

describe("FapeGalleryComponent paginator fetch", () => {
  let fixture: ComponentFixture<FapeGalleryComponent>;
  let httpMock: HttpTestingController;

  function configure(isMobile: boolean): void {
    TestBed.configureTestingModule({
      declarations: [
        FapeGalleryComponent,
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
        { provide: MatDialog, useValue: {} },
        { provide: ConfirmDialog, useValue: { show: () => undefined } },
        { provide: I18n, useValue: { trl: (_m: string, k: string) => k } },
        { provide: NavigationService, useValue: {} },
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
    fixture = TestBed.createComponent(FapeGalleryComponent);
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
    fixture = TestBed.createComponent(FapeGalleryComponent);
    fixture.detectChanges();

    const reqs = httpMock.match(LIST_URL);
    reqs.forEach((r) =>
      r.flush({ data: { payload: [], paging: { total: 0 } } })
    );

    expect(reqs.length).toBe(1);
    expect(reqs[0].request.body.paging.limit).toBe(10);
  });
});
