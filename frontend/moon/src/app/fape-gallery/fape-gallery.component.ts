import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ConfirmDialog } from 'src/common/dialog';
import { ControlledPaginatorComponent } from '../controlled-paginator/controlled-paginator.component';
import { DirTreeNavComponent } from '../dir-tree-nav/dir-tree-nav.component';
import { I18n } from '../i18n.service';
import { NavigationService } from '../navigation.service';
import { NavType } from '../routes';

interface CreateFapeGalleryReq {
  name: string;
  dirFileKey: string;
  maxPages?: number;
}

interface ListFapeGalleryReq {
  paging: { page: number; limit: number };
}

interface ListedFapeGallery {
  galleryNo: string;
  name: string;
  dirFileKey: string;
  enabled: boolean;
  maxPages: number;
  createdAt: string;
}

interface ToggleFapeGalleryReq {
  galleryNo: string;
  enabled: boolean;
}

interface TriggerFapeGalleryReq {
  galleryNo: string;
  maxPage: number;
}

interface DeleteFapeGalleryReq {
  galleryNo: string;
}

@Component({
  selector: 'app-fape-gallery',
  templateUrl: './fape-gallery.component.html',
  styleUrls: ['./fape-gallery.component.css'],
})
export class FapeGalleryComponent implements OnInit {
  headers: string[] = ['name', 'dirFileKey', 'enabled', 'maxPages', 'createdAt', 'operation'];

  // Create gallery
  showCreatePanel = false;
  createGalleryReq: CreateFapeGalleryReq = {
    name: '',
    dirFileKey: '',
    maxPages: 3,
  };
  createDirName = '';
  prevSelectedFileKey: string | undefined;

  // Gallery list
  galleries: ListedFapeGallery[] = [];
  listGalleryReq: ListFapeGalleryReq = {
    paging: { page: 1, limit: 10 },
  };
  total = 0;

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  @ViewChild(DirTreeNavComponent)
  dirTreeNav: DirTreeNavComponent;

  constructor(
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private dialog: MatDialog,
    private confirmDialog: ConfirmDialog,
    private i18n: I18n,
    private nav: NavigationService
  ) {}

  ngOnInit(): void {
    this.listGalleries();
  }

  listGalleries() {
    this.listGalleryReq.paging = this.pagingController
      ? this.pagingController.paging
      : { page: 1, limit: 10 };
    this.http
      .post<any>(`/drone/open/api/fape/gallery/list`, this.listGalleryReq)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
            return;
          }
          const dat = resp.data;
          this.galleries = dat.payload || [];
          this.total = dat.paging?.total || 0;
          if (this.pagingController) {
            this.pagingController.onTotalChanged(dat.paging);
          }
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open(
            this.i18n.trl('fape-gallery', 'requestFailedUnknownError'),
            'ok',
            { duration: 3000 }
          );
        },
      });
  }

  onPageChanged() {
    this.listGalleries();
  }

  toggleCreatePanel(): void {
    this.showCreatePanel = !this.showCreatePanel;
    if (this.showCreatePanel) {
      this.prevSelectedFileKey = undefined;
    }
  }

  selectDir(n) {
    this.createGalleryReq.dirFileKey = n.fileKey;
    this.createDirName = n.name;
    this.dirTreeNav.collapseAll();
  }

  createGallery() {
    if (!this.createGalleryReq.name.trim() || !this.createGalleryReq.dirFileKey) {
      this.snackBar.open(
        this.i18n.trl('fape-gallery', 'nameAndDirRequired'),
        'ok',
        { duration: 3000 }
      );
      return;
    }
    const req: CreateFapeGalleryReq = {
      name: this.createGalleryReq.name.trim(),
      dirFileKey: this.createGalleryReq.dirFileKey,
      maxPages: this.createGalleryReq.maxPages || 3,
    };
    this.http
      .post<any>(`/drone/open/api/fape/gallery/create`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
            return;
          }
          this.snackBar.open(
            `${this.i18n.trl('fape-gallery', 'galleryCreated')}: ${resp.data.galleryNo}`,
            'ok',
            { duration: 5000 }
          );
          this.createGalleryReq = { name: '', dirFileKey: '', maxPages: 3 };
          this.createDirName = '';
          this.showCreatePanel = false;
          this.listGalleries();
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open(
            this.i18n.trl('fape-gallery', 'requestFailedUnknownError'),
            'ok',
            { duration: 3000 }
          );
        },
      });
  }

  toggleGalleryEnabled(g: ListedFapeGallery): void {
    const newEnabled = !g.enabled;
    const req: ToggleFapeGalleryReq = {
      galleryNo: g.galleryNo,
      enabled: newEnabled,
    };
    this.http.post<any>('/drone/open/api/fape/gallery/toggle', req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
          return;
        }
        g.enabled = newEnabled;
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open(
          this.i18n.trl('fape-gallery', 'requestFailedUnknownError'),
          'ok',
          { duration: 3000 }
        );
      },
    });
  }

  clickTrigger(g: ListedFapeGallery) {
    this.dialog
      .open(TriggerFapeGalleryDialogComponent, {
        width: '450px',
        data: {
          galleryNo: g.galleryNo,
          maxPage: 3,
        },
      })
      .afterClosed()
      .subscribe((result) => {
        if (result) {
          let maxPage: number;
          if (result.maxPage === '' || result.maxPage === null || result.maxPage === undefined) {
            maxPage = 3; // default: 3 pages
          } else {
            maxPage = Number(result.maxPage);
            if (isNaN(maxPage) || maxPage < 0) {
              maxPage = 0; // <= 0 = full gallery
            }
          }
          this.triggerGallery(g.galleryNo, maxPage);
        }
      });
  }

  triggerGallery(galleryNo: string, maxPage: number) {
    const req: TriggerFapeGalleryReq = {
      galleryNo,
      maxPage,
    };
    this.http
      .post<any>(`/drone/open/api/fape/gallery/trigger`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
            return;
          }
          this.snackBar.open(
            this.i18n.trl('fape-gallery', 'triggered'),
            'ok',
            { duration: 3000 }
          );
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open(
            this.i18n.trl('fape-gallery', 'requestFailedUnknownError'),
            'ok',
            { duration: 3000 }
          );
        },
      });
  }

  gotoDir(g: ListedFapeGallery) {
    this.nav.navigateTo(NavType.MANAGE_FILES, [{ parentDirKey: g.dirFileKey }]);
  }

  deleteGallery(g: ListedFapeGallery) {
    this.confirmDialog.show(
      this.i18n.trl('fape-gallery', 'deleteGallery'),
      [`${this.i18n.trl('fape-gallery', 'sureToDelete')} '${g.name}'`],
      () => {
        const req: DeleteFapeGalleryReq = { galleryNo: g.galleryNo };
        this.http
          .post<any>(`/drone/open/api/fape/gallery/delete`, req)
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
                return;
              }
              this.listGalleries();
            },
            error: (err) => {
              console.log(err);
              this.snackBar.open(
                this.i18n.trl('fape-gallery', 'requestFailedUnknownError'),
                'ok',
                { duration: 3000 }
              );
            },
          });
      }
    );
  }
}

@Component({
  selector: 'app-trigger-fape-gallery',
  template: `
    <h1 mat-dialog-title>{{ 'fape-gallery' | trl:'triggerGallery' }}</h1>
    <div mat-dialog-content>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'fape-gallery' | trl:'galleryNo' }}</mat-label>
        <input readonly disabled matInput [ngModel]="data.galleryNo" />
      </mat-form-field>
      <mat-form-field style="width: 400px">
        <mat-label>{{ 'fape-gallery' | trl:'maxPageOptional' }}</mat-label>
        <input matInput type="number" [(ngModel)]="data.maxPage" />
      </mat-form-field>
    </div>
    <div mat-dialog-actions class="d-flex justify-content-end">
      <button mat-button [mat-dialog-close]="{maxPage: data.maxPage}">
        {{ 'fape-gallery' | trl:'trigger' }}
      </button>
      <button mat-button [mat-dialog-close]="null" cdkFocusInitial>
        {{ 'fape-gallery' | trl:'no' }}
      </button>
    </div>
  `,
})
export class TriggerFapeGalleryDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<TriggerFapeGalleryDialogComponent, any>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}
}
