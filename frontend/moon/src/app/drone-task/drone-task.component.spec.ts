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
import { By } from "@angular/platform-browser";
import { MatTreeModule } from "@angular/material/tree";

import {
  FakeMediaQueryList,
  installFakeMatchMedia,
} from "src/testing/fake-media-query-list";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { DirTopDownTreeNode } from "src/common/dir-tree";
import { DirTreeNavComponent } from "../dir-tree-nav/dir-tree-nav.component";
import { I18n } from "../i18n.service";
import { NavigationService } from "../navigation.service";
import { TrlPipe } from "../trl.pipe";
import { DroneTaskComponent } from "./drone-task.component";

const LIST_URL = "/drone/open/api/list-task";

describe("DroneTaskComponent search", () => {
  let fixture: ComponentFixture<DroneTaskComponent>;
  let httpMock: HttpTestingController;
  let fake: FakeMediaQueryList;

  beforeEach(() => {
    fake = installFakeMatchMedia(false);
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

  it("shows the wide columns at 1280px, the compact columns at 390px, and follows a flip", () => {
    const component = fixture.componentInstance;

    const wide = [
      "thumbnail",
      "taskId",
      "status",
      "url",
      "platform",
      "attempt",
      "lastPage",
      "dirName",
      "fileCount",
      "updatedAt",
      "remark",
      "operation",
    ];
    const compact = ["thumbnail", "status", "dirName", "operation"];

    expect(component.displayedColumns).toEqual(wide);

    fake.fireChange(true);
    expect(component.displayedColumns).toEqual(compact);

    fake.fireChange(false);
    expect(component.displayedColumns).toEqual(wide);
  });
});

describe("DroneTaskComponent directory tree", () => {
  const TREE_URL = "/vfm/open/api/file/dir/top-down-tree";
  const LAST_DIR_URL = "/drone/open/api/task/last-selected-dir";

  let fixture: ComponentFixture<DroneTaskComponent>;
  let httpMock: HttpTestingController;
  let root: DirTopDownTreeNode;
  let groupDocs: DirTopDownTreeNode;
  let nodeDeep: DirTopDownTreeNode;
  let nodeOther: DirTopDownTreeNode;

  beforeEach(() => {
    installFakeMatchMedia(false);
    TestBed.configureTestingModule({
      declarations: [
        DroneTaskComponent,
        ControlledPaginatorComponent,
        DirTreeNavComponent,
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
        MatTreeModule,
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
      ],
    });
    httpMock = TestBed.inject(HttpTestingController);
    fixture = TestBed.createComponent(DroneTaskComponent);
    fixture.detectChanges();

    nodeDeep = { fileKey: "deep", name: "deep" };
    groupDocs = { name: "docs", child: [nodeDeep] };
    nodeOther = { fileKey: "other", name: "other" };
    root = { fileKey: "root", name: "root", child: [groupDocs, nodeOther] };
  });

  afterEach(() => {
    httpMock.verify();
  });

  function openPanelAndLoadTree(): DirTreeNavComponent {
    const component = fixture.componentInstance;
    component.toggleCreateTaskPanelShown();
    fixture.detectChanges();

    httpMock
      .match(LIST_URL)
      .forEach((r) =>
        r.flush({ data: { payload: [], paging: { total: 0 } } })
      );
    httpMock
      .match("/drone/open/api/list-platforms")
      .forEach((r) => r.flush({ data: [] }));
    httpMock
      .match("/drone/open/api/list-last-page-platforms")
      .forEach((r) => r.flush({ data: [] }));

    httpMock.expectOne(TREE_URL).flush({ data: root });
    fixture.detectChanges();

    return fixture.debugElement.query(By.directive(DirTreeNavComponent))
      .componentInstance;
  }

  /**
   * Answers the remembered-directory request, then lets the tree apply the key
   * the way the template binding does: the deferred emit lands after the
   * change-detection pass, so the page updates without a dev-mode mismatch.
   */
  async function applyRememberedDir(fileKey: string): Promise<void> {
    httpMock.expectOne(LAST_DIR_URL).flush({ data: { fileKey } });
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();
    await fixture.whenStable();
  }

  function parentDirectoryValue(): string {
    const fields = Array.from(
      fixture.nativeElement.querySelectorAll("mat-form-field")
    ) as HTMLElement[];
    const field = fields.find(
      (f) =>
        (f.querySelector("mat-label")?.textContent || "").trim() ===
        "parentDirectory"
    );
    if (!field) {
      throw new Error("no parent directory field");
    }
    return (field.querySelector("input") as HTMLInputElement).value;
  }

  function treeLabel(text: string): HTMLElement {
    const labels = Array.from(
      fixture.nativeElement.querySelectorAll(".dir-tree-name")
    ) as HTMLElement[];
    const label = labels.find((l) => (l.textContent || "").trim() === text);
    if (!label) {
      throw new Error(`no tree label reading '${text}'`);
    }
    return label;
  }

  function selectedTreeLabel(): HTMLElement {
    const labels = Array.from(
      fixture.nativeElement.querySelectorAll(".dir-tree-name-selected")
    ) as HTMLElement[];
    expect(labels.length).toBe(1);
    return labels[0];
  }

  /** Empties the tree's search box the way a user would. */
  function clearTreeSearch(tree: DirTreeNavComponent): void {
    expect(tree.searchDirTreeName).toBeTruthy();
    const input = fixture.debugElement.query(By.directive(DirTreeNavComponent))
      .nativeElement.querySelector("input") as HTMLInputElement;
    input.value = "";
    input.dispatchEvent(new Event("input"));
    input.dispatchEvent(new KeyboardEvent("keyup", { key: "Backspace" }));
    fixture.detectChanges();
  }

  it("applies the last selected directory to the tree and the parent directory field", async () => {
    const tree = openPanelAndLoadTree();

    await applyRememberedDir("deep");

    expect(fixture.componentInstance.createTaskDirName).toBe("deep");
    expect(parentDirectoryValue()).toBe("deep");

    const selected = selectedTreeLabel();
    expect(selected.textContent).toContain("/deep");
    expect(selected.getAttribute("aria-selected")).toBe("true");
    expect(tree.dirTreeControl.isExpanded(root)).toBe(true);
    expect(tree.dirTreeControl.isExpanded(groupDocs)).toBe(true);
    expect(tree.searchDirTreeName).toBe("deep");

    const rendered = Array.from(
      fixture.nativeElement.querySelectorAll(".dir-tree-name")
    ).map((l: HTMLElement) => (l.textContent || "").trim());
    expect(rendered).not.toContain("/other");
  });

  it("leaves the tree expanded after a click and rebinds the parent directory", async () => {
    const tree = openPanelAndLoadTree();
    await applyRememberedDir("deep");
    clearTreeSearch(tree);

    treeLabel("/other").click();
    fixture.detectChanges();
    await fixture.whenStable();

    expect(fixture.componentInstance.createTaskDirName).toBe("other");
    expect(fixture.componentInstance.createTaskReq.dirFileKey).toBe("other");
    expect(parentDirectoryValue()).toBe("other");

    const selected = selectedTreeLabel();
    expect(selected.textContent).toContain("/other");
    expect(selected.getAttribute("aria-selected")).toBe("true");
    expect(tree.dirTreeControl.isExpanded(root)).toBe(true);
    expect(tree.dirTreeControl.isExpanded(groupDocs)).toBe(true);
  });
});
