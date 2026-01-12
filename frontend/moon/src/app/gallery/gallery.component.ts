import { Component, OnInit, ViewChild } from "@angular/core";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { MatPaginator } from "@angular/material/paginator";
import {
  animateElementExpanding,
  getExpanded,
  isIdEqual,
} from "src/animate/animate-util";
import { Gallery } from "src/common/gallery";
import { Paging } from "src/common/paging";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { NavigationService } from "../navigation.service";
import { NavType } from "../routes";
import { Env } from "src/common/env-util";
import { GalleryAccessComponent } from "../gallery-access/gallery-access.component";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { I18n } from "../i18n.service";

export interface ListGalleriesResp {
  paging: Paging;
  payload: Gallery[];
}

@Component({
  selector: "app-gallery",
  templateUrl: "./gallery.component.html",
  styleUrls: ["./gallery.component.css"],
  animations: [animateElementExpanding()],
})
export class GalleryComponent implements OnInit {
  readonly DESKTOP_COLUMNS = [
    "thumbnail",
    "galleryNo",
    "name",
    "createTime",
    "updateTime",
    "createBy",
  ];
  readonly MOBILE_COLUMNS = ["thumbnail", "name"];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  galleries: Gallery[] = [];
  expandedElement: Gallery = null;
  newGalleryName: string = "";
  showCreateGalleryDiv: boolean = false;

  trl = (k) => {
    return this.i18n.trl("gallery", k);
  };
  idEquals = isIdEqual;
  getExpandedEle = (row) => getExpanded(row, this.expandedElement);

  constructor(
    public env: Env,
    private http: HttpClient,
    private navigation: NavigationService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar,
    public i18n: I18n
  ) {}

  ngOnInit() {}

  fetchGalleries() {
    this.http
      .post<any>(`vfm/open/api/gallery/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.pagingController.onTotalChanged(resp.data.paging);
          this.galleries = resp.data.payload;
          for (let g of this.galleries) {
            if (g.thumbnailToken) {
              g.thumbnailUrl =
                "fstore/file/raw?key=" + encodeURIComponent(g.thumbnailToken);
            }
          }
          this.expandedElement = null;
        },
      });
  }

  createGallery() {
    if (!this.newGalleryName) {
      this.snackBar.open("Please enter new gallery's name", "ok", {
        duration: 3000,
      });
      return;
    }

    this.http
      .post<any>(`vfm/open/api/gallery/new`, {
        name: this.newGalleryName,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.newGalleryName = null;
        },
        complete: () => {
          this.fetchGalleries();
          this.expandedElement = null;
          this.showCreateGalleryDiv = false;
        },
      });
  }

  // todo (impl this later)
  shareGallery(g: Gallery) {}

  deleteGallery(galleryNo: string, galleryName: string) {
    if (!galleryNo) return;

    this.dialog
      .open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Delete Gallery",
          msg: [`You sure you want to delete '${galleryName}'`],
          isNoBtnDisplayed: true,
        },
      })
      .afterClosed()
      .subscribe((confirm) => {
        if (!confirm) {
          this.expandedElement = null;
          return;
        }

        this.http
          .post<any>(`vfm/open/api/gallery/delete`, {
            galleryNo: galleryNo,
          })
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
            },
            complete: () => this.fetchGalleries(),
          });
        this.expandedElement = null;
      });
  }

  browse(galleryNo: string) {
    this.navigation.navigateTo(NavType.GALLERY_IMAGE, [
      { galleryNo: galleryNo },
    ]);
  }

  updateGallery(galleryNo: string, name: string) {
    if (!galleryNo || !name) return;

    this.http
      .post<any>(`vfm/open/api/gallery/update`, {
        galleryNo: galleryNo,
        name: name,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
        },
        complete: () => {
          this.expandedElement = null;
          this.fetchGalleries();
        },
      });
  }

  popToGrantAccess(g: Gallery): void {
    if (!g) return;

    const dialogRef: MatDialogRef<GalleryAccessComponent, boolean> =
      this.dialog.open(GalleryAccessComponent, {
        width: "700px",
        data: { galleryNo: g.galleryNo, name: g.name },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      // do nothing
    });
  }

  openFileDir(g: Gallery) {
    if (!g.dirFileKey || !g.isOwner) {
      return;
    }
    this.navigation.navigateTo(NavType.MANAGE_FILES, [
      { parentDirKey: g.dirFileKey },
    ]);
  }
}
