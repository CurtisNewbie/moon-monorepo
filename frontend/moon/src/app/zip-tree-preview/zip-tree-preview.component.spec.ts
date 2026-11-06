import { CommonModule } from '@angular/common';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatTreeModule } from '@angular/material/tree';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import {
  ZipTreePreview,
  ZipTreeNode,
  ZIP_TREE_TRUNCATED_BY_DEPTH,
  ZIP_TREE_TRUNCATED_BY_ENTRIES,
} from 'src/common/zip-tree';
import { I18n } from '../i18n.service';
import { TrlPipe } from '../trl.pipe';
import {
  ZipTreePreviewComponent,
  ZipTreePreviewData,
} from './zip-tree-preview.component';

function folder(name: string, children: ZipTreeNode[] = []): ZipTreeNode {
  return { name, dir: true, children };
}

function file(name: string, size: number): ZipTreeNode {
  return { name, dir: false, size };
}

function preview(
  children: ZipTreeNode[],
  extra: Partial<ZipTreePreview> = {}
): ZipTreePreview {
  return {
    tree: folder('', children),
    truncated: false,
    truncatedReasons: [],
    includedEntries: 0,
    ...extra,
  };
}

describe('ZipTreePreviewComponent', () => {
  let fixture: ComponentFixture<ZipTreePreviewComponent>;
  let component: ZipTreePreviewComponent;
  let dialogData: ZipTreePreviewData;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ZipTreePreviewComponent, TrlPipe],
      imports: [
        CommonModule,
        NoopAnimationsModule,
        MatButtonModule,
        MatTreeModule,
      ],
      providers: [
        { provide: MatDialogRef, useValue: { close: () => undefined } },
        { provide: MAT_DIALOG_DATA, useFactory: () => dialogData },
        {
          provide: I18n,
          useValue: {
            trl: (mod: string, key: string, ...args: any[]) =>
              args.length ? `${key} ${args.join(' ')}` : key,
          },
        },
      ],
    });
  });

  function open(data: ZipTreePreviewData): void {
    dialogData = data;
    fixture = TestBed.createComponent(ZipTreePreviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  /** The rows the user can see: a collapsed folder keeps its subtree in the DOM. */
  function rows(): HTMLElement[] {
    return (
      Array.from(
        fixture.nativeElement.querySelectorAll('.zip-tree-row')
      ) as HTMLElement[]
    ).filter((row) => row.getClientRects().length > 0);
  }

  function rowTexts(): string[] {
    return rows().map((row) =>
      (row.querySelector('.zip-tree-name') as HTMLElement).textContent.trim()
    );
  }

  function rowNamed(name: string): HTMLElement {
    return rows().find(
      (row) =>
        (row.querySelector('.zip-tree-name') as HTMLElement).textContent.trim() ===
        name
    );
  }

  function toggleNamed(name: string): HTMLButtonElement {
    return rowNamed(name).querySelector(
      'button.zip-tree-toggle'
    ) as HTMLButtonElement;
  }

  it('shows each top-level folder and file by name', () => {
    open({
      name: 'photos.zip',
      preview: preview([folder('album'), file('readme.txt', 2048)]),
    });

    expect(rowTexts()).toEqual(['album', 'readme.txt']);
  });

  it('shows a file size with the app size formatter and no size on a folder', () => {
    open({
      name: 'photos.zip',
      preview: preview([folder('album'), file('readme.txt', 2048)]),
    });

    const fileRow = rowNamed('readme.txt');
    expect(
      (fileRow.querySelector('.zip-tree-size') as HTMLElement).textContent.trim()
    ).toBe('2.0 kb');
    expect(rowNamed('album').querySelector('.zip-tree-size')).toBeNull();
  });

  it('shows a missing size as zero', () => {
    // a zero-byte file omits size from the payload
    open({
      name: 'photos.zip',
      preview: preview([{ name: 'empty.txt', dir: false }]),
    });

    expect(
      (rowNamed('empty.txt').querySelector('.zip-tree-size') as HTMLElement)
        .textContent.trim()
    ).toBe('0.0 kb');
  });

  it('marks folders and files with the app icons', () => {
    open({
      name: 'photos.zip',
      preview: preview([folder('album'), file('readme.txt', 10)]),
    });

    expect(
      rowNamed('album').querySelector('i.bi-folder-fill, i.bi-folder2-open')
    ).not.toBeNull();
    expect(rowNamed('readme.txt').querySelector('i.bi-file-text')).not.toBeNull();
  });

  it('starts top-level folders expanded and deeper folders collapsed', () => {
    open({
      name: 'photos.zip',
      preview: preview([
        folder('top', [folder('deeper', [file('hidden.txt', 1)]), file('shown.txt', 2)]),
      ]),
    });

    expect(rowTexts()).toEqual(['top', 'deeper', 'shown.txt']);
  });

  it('opens a deeper folder on its toggle', () => {
    open({
      name: 'photos.zip',
      preview: preview([
        folder('top', [folder('deeper', [file('hidden.txt', 1)]), file('shown.txt', 2)]),
      ]),
    });

    toggleNamed('deeper').click();
    fixture.detectChanges();

    expect(rowTexts()).toEqual(['top', 'deeper', 'hidden.txt', 'shown.txt']);
  });

  it('names the cause of a truncated tree', () => {
    open({
      name: 'huge.zip',
      preview: preview([file('a.bin', 1)], {
        truncated: true,
        truncatedReasons: [ZIP_TREE_TRUNCATED_BY_DEPTH, ZIP_TREE_TRUNCATED_BY_ENTRIES],
        includedEntries: 1,
      }),
    });

    const notice = fixture.nativeElement.querySelector(
      '.zip-tree-truncated'
    ) as HTMLElement;
    expect(notice).not.toBeNull();
    expect(notice.textContent).toContain('truncatedTitle');
    expect(notice.textContent).toContain('truncatedDepth');
    expect(notice.textContent).toContain('truncatedEntries');
  });

  it('names only the cause the server reported', () => {
    open({
      name: 'deep.zip',
      preview: preview([file('a.bin', 1)], {
        truncated: true,
        truncatedReasons: [ZIP_TREE_TRUNCATED_BY_DEPTH],
        includedEntries: 1,
      }),
    });

    const notice = fixture.nativeElement.querySelector(
      '.zip-tree-truncated'
    ) as HTMLElement;
    expect(notice.textContent).toContain('truncatedDepth');
    expect(notice.textContent).not.toContain('truncatedEntries');
  });

  it('shows no cut notice for a complete tree', () => {
    open({
      name: 'small.zip',
      preview: preview([file('a.bin', 1)], { includedEntries: 1 }),
    });

    expect(fixture.nativeElement.querySelector('.zip-tree-truncated')).toBeNull();
  });

  it('shows an empty state for an archive with no entries', () => {
    open({ name: 'empty.zip', preview: preview([]) });

    const empty = fixture.nativeElement.querySelector(
      '.zip-tree-empty'
    ) as HTMLElement;
    expect(empty.textContent).toContain('emptyTitle');
    expect(empty.textContent).toContain('emptyHint');
    expect(fixture.nativeElement.querySelector('mat-tree')).toBeNull();
  });

  it('renders a crafted entry name as text, never as markup', () => {
    const crafted = '<img src="x" onerror="alert(1)">.txt';
    open({ name: 'evil.zip', preview: preview([file(crafted, 1)]) });

    expect(rowNamed(crafted)).toBeDefined();
    expect(
      (rowNamed(crafted).querySelector('.zip-tree-name') as HTMLElement)
        .querySelector('img')
    ).toBeNull();
  });

  it('titles the dialog with the archive name and counts the entries shown', () => {
    open({
      name: 'photos.zip',
      preview: preview([file('readme.txt', 2048)], { includedEntries: 1 }),
    });

    const title = fixture.nativeElement.querySelector(
      '.zip-tree-title'
    ) as HTMLElement;
    expect(title.textContent).toContain('photos.zip');
    expect(
      (fixture.nativeElement.querySelector('.zip-tree-meta') as HTMLElement)
        .textContent
    ).toContain('count 1');
  });
});
