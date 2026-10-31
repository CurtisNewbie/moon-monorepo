import { CommonModule } from '@angular/common';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTreeModule } from '@angular/material/tree';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { DirTopDownTreeNode } from 'src/common/dir-tree';
import { I18n } from '../i18n.service';
import { NavigationService } from '../navigation.service';
import { TrlPipe } from '../trl.pipe';
import { DirTreeNavComponent } from './dir-tree-nav.component';

const TREE_URL = '/vfm/open/api/file/dir/top-down-tree';

describe('DirTreeNavComponent selection', () => {
  let fixture: ComponentFixture<DirTreeNavComponent>;
  let component: DirTreeNavComponent;
  let httpMock: HttpTestingController;
  let emitted: DirTopDownTreeNode[];

  // root > a (grouping node, no file key) > b > c, plus a sibling x.
  let root: DirTopDownTreeNode;
  let groupA: DirTopDownTreeNode;
  let nodeB: DirTopDownTreeNode;
  let nodeC: DirTopDownTreeNode;
  let nodeX: DirTopDownTreeNode;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [DirTreeNavComponent, TrlPipe],
      imports: [
        CommonModule,
        HttpClientTestingModule,
        NoopAnimationsModule,
        FormsModule,
        MatButtonModule,
        MatTreeModule,
        MatFormFieldModule,
        MatInputModule,
      ],
      providers: [
        { provide: MatSnackBar, useValue: { open: () => undefined } },
        {
          provide: I18n,
          useValue: { trl: (moduleName: string, key: string) => key },
        },
        { provide: NavigationService, useValue: {} },
      ],
    });

    fixture = TestBed.createComponent(DirTreeNavComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
    emitted = [];
    component.selectedEmiter.subscribe((node) => emitted.push(node));

    nodeC = { fileKey: 'c', name: 'c' };
    nodeB = { fileKey: 'b', name: 'b', child: [nodeC] };
    groupA = { name: 'a', child: [nodeB] };
    nodeX = { fileKey: 'x', name: 'x' };
    root = { fileKey: 'root', name: 'root', child: [groupA, nodeX] };
  });

  afterEach(() => {
    httpMock.verify();
  });

  function loadTree(): void {
    fixture.detectChanges(); // starts ngOnInit and the tree fetch
    httpMock.expectOne(TREE_URL).flush({ data: root });
    fixture.detectChanges();
  }

  function selectedLabels(): HTMLElement[] {
    return Array.from(
      fixture.nativeElement.querySelectorAll('.dir-tree-name-selected')
    ) as HTMLElement[];
  }

  function labelText(text: string): HTMLElement {
    const labels = Array.from(
      fixture.nativeElement.querySelectorAll('.dir-tree-name')
    ) as HTMLElement[];
    const label = labels.find((l) => (l.textContent || '').trim() === text);
    if (!label) {
      throw new Error(`no tree label reading '${text}'`);
    }
    return label;
  }

  it('applies a key set before the tree arrives once the tree has loaded', async () => {
    component.selectedFileKey = 'c';

    loadTree();
    await fixture.whenStable();

    expect(emitted.length).toBe(1);
    expect(emitted[0].fileKey).toBe('c');
  });

  it('emits a shallow copy of the selected node', async () => {
    loadTree();

    component.selectedFileKey = 'b';
    fixture.detectChanges();
    await fixture.whenStable();

    expect(emitted.length).toBe(1);
    expect(emitted[0]).toEqual(nodeB);
    expect(emitted[0]).not.toBe(nodeB);
  });

  it('emits after the change-detection pass that applies the input', async () => {
    loadTree();

    component.selectedFileKey = 'b';
    fixture.detectChanges();
    expect(emitted.length).toBe(0);

    await fixture.whenStable();

    expect(emitted.length).toBe(1);
    expect(emitted[0].fileKey).toBe('b');
  });

  it('expands every ancestor on the path, including a grouping node without a key', () => {
    loadTree();

    component.selectedFileKey = 'c';
    fixture.detectChanges();

    expect(component.dirTreeControl.isExpanded(root)).toBe(true);
    expect(component.dirTreeControl.isExpanded(groupA)).toBe(true);
    expect(component.dirTreeControl.isExpanded(nodeB)).toBe(true);
    expect(component.dirTreeControl.isExpanded(nodeC)).toBe(false);
    expect(component.dirTreeControl.isExpanded(nodeX)).toBe(false);
  });

  it('marks exactly the selected node with the class and aria-selected', () => {
    loadTree();

    component.selectedFileKey = 'c';
    fixture.detectChanges();

    const labels = selectedLabels();
    expect(labels.length).toBe(1);
    expect(labels[0].textContent).toContain('/c');
    expect(labels[0].getAttribute('aria-selected')).toBe('true');

    const marked = fixture.nativeElement.querySelectorAll(
      '[aria-selected="true"]'
    );
    expect(marked.length).toBe(1);
  });

  it('scrolls only the tree box so the selected label is centred', () => {
    const leaves: DirTopDownTreeNode[] = [];
    for (let i = 0; i < 41; i++) {
      leaves.push({ fileKey: `k${i}`, name: `k${i}` });
    }
    fixture.detectChanges();
    httpMock
      .expectOne(TREE_URL)
      .flush({ data: { fileKey: 'root', name: 'root', child: leaves } });
    fixture.detectChanges();

    const scrollIntoView = spyOn(Element.prototype, 'scrollIntoView');

    component.selectedFileKey = 'k20';
    fixture.detectChanges();

    expect(scrollIntoView).not.toHaveBeenCalled();

    const box = fixture.nativeElement.querySelector(
      '.dir-tree-box'
    ) as HTMLElement;
    expect(box).toBeTruthy();
    expect(box.scrollTop).toBeGreaterThan(0);

    const label = selectedLabels()[0];
    const boxRect = box.getBoundingClientRect();
    const labelRect = label.getBoundingClientRect();
    const boxCentre = boxRect.top + boxRect.height / 2;
    const labelCentre = labelRect.top + labelRect.height / 2;
    expect(Math.abs(labelCentre - boxCentre)).toBeLessThan(2);
  });

  it('applies a key set after the tree has loaded immediately', async () => {
    loadTree();

    component.selectedFileKey = 'b';
    fixture.detectChanges();
    await fixture.whenStable();

    expect(emitted.length).toBe(1);
    expect(component.dirTreeControl.isExpanded(root)).toBe(true);
    expect(selectedLabels().length).toBe(1);
  });

  it('emits once when the same key is applied twice', async () => {
    loadTree();

    component.selectedFileKey = 'b';
    component.selectedFileKey = 'b';
    fixture.detectChanges();
    await fixture.whenStable();

    expect(emitted.length).toBe(1);
  });

  it('marks a clicked node, emits once, and collapses nothing', () => {
    loadTree();
    component.dirTreeControl.expand(root);
    component.dirTreeControl.expand(groupA);
    component.dirTreeControl.expand(nodeB);
    fixture.detectChanges();

    labelText('/b').click();
    fixture.detectChanges();

    expect(emitted.length).toBe(1);
    expect(emitted[0].fileKey).toBe('b');
    expect(selectedLabels().length).toBe(1);
    expect(selectedLabels()[0].textContent).toContain('/b');
    expect(component.dirTreeControl.isExpanded(root)).toBe(true);
    expect(component.dirTreeControl.isExpanded(groupA)).toBe(true);
    expect(component.dirTreeControl.isExpanded(nodeB)).toBe(true);
  });

  it('emits again when the already selected node is clicked', () => {
    loadTree();
    component.dirTreeControl.expand(root);
    component.dirTreeControl.expand(groupA);
    component.dirTreeControl.expand(nodeB);
    fixture.detectChanges();

    labelText('/b').click();
    fixture.detectChanges();
    expect(emitted.length).toBe(1);

    labelText('/b').click();
    fixture.detectChanges();

    expect(emitted.map((n) => n.fileKey)).toEqual(['b', 'b']);
  });

  it('highlights with the theme pink at 0.16 alpha', () => {
    loadTree();

    component.selectedFileKey = 'c';
    fixture.detectChanges();

    expect(getComputedStyle(selectedLabels()[0]).backgroundColor).toBe(
      'rgba(233, 30, 99, 0.16)'
    );
  });

  it('clears the highlight and emits nothing for a key absent from the tree', async () => {
    loadTree();
    component.selectedFileKey = 'c';
    fixture.detectChanges();
    await fixture.whenStable();
    expect(selectedLabels().length).toBe(1);
    expect(emitted.length).toBe(1);

    component.selectedFileKey = 'missing';
    fixture.detectChanges();
    await fixture.whenStable();

    expect(selectedLabels().length).toBe(0);
    expect(emitted.length).toBe(1);
  });

  it('keeps the selected node highlighted while a search shows results', () => {
    loadTree();
    component.selectedFileKey = 'c';
    fixture.detectChanges();

    component.searchDirTreeName = 'c';
    component.onSearchDirTreeNameChanged();
    fixture.detectChanges();

    expect(component.dirTreeControl.isExpanded(nodeB)).toBe(true);
    expect(labelText('/c').classList).toContain('dir-tree-name-selected');
    expect(selectedLabels().length).toBe(1);
  });

  it('restores the selected path and reveals it when the search is cleared', () => {
    loadTree();
    component.selectedFileKey = 'c';
    fixture.detectChanges();

    component.searchDirTreeName = 'zzz';
    component.onSearchDirTreeNameChanged();
    fixture.detectChanges();
    expect(component.dirTreeControl.isExpanded(root)).toBe(false);

    const scrollSpy = spyOn(Element.prototype, 'scrollIntoView');

    component.searchDirTreeName = '';
    component.onSearchDirTreeNameChanged();
    fixture.detectChanges();

    expect(component.dirTreeControl.isExpanded(root)).toBe(true);
    expect(component.dirTreeControl.isExpanded(groupA)).toBe(true);
    expect(component.dirTreeControl.isExpanded(nodeB)).toBe(true);
    expect(selectedLabels().length).toBe(1);
    expect(selectedLabels()[0].textContent).toContain('/c');
    expect(scrollSpy).not.toHaveBeenCalled();
  });

  it('emits nothing and refetches nothing when the search is cleared', async () => {
    loadTree();
    component.selectedFileKey = 'c';
    fixture.detectChanges();
    await fixture.whenStable();
    expect(emitted.length).toBe(1);

    component.searchDirTreeName = 'c';
    component.onSearchDirTreeNameChanged();
    fixture.detectChanges();

    component.searchDirTreeName = '';
    component.onSearchDirTreeNameChanged();
    fixture.detectChanges();
    await fixture.whenStable();

    expect(emitted.length).toBe(1);
    httpMock.expectNone(TREE_URL);
  });

  it('drops a pending input emit when a click selects another node first', async () => {
    loadTree();
    component.dirTreeControl.expand(root);
    fixture.detectChanges();

    component.selectedFileKey = 'c';
    fixture.detectChanges();
    labelText('/x').click();

    expect(emitted.length).toBe(1);
    expect(emitted[0].fileKey).toBe('x');

    await fixture.whenStable();

    expect(emitted.map((n) => n.fileKey)).toEqual(['x']);
  });
});
