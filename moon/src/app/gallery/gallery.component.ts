import { Component, OnInit, ViewChild } from "@angular/core";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { MatPaginator } from "@angular/material/paginator";
import {
  animateElementExpanding,
  getExpanded,
  isIdEqual,
} from "src/animate/animate-util";
import { environment } from "src/environments/environment";
import { Gallery } from "src/common/gallery";
import { Paging, PagingController } from "src/common/paging";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { NavigationService } from "../navigation.service";
import { Toaster } from "../notification.service";
import { NavType } from "../routes";
import { isMobile } from "src/common/env-util";
import { GalleryAccessComponent } from "../gallery-access/gallery-access.component";
import { HttpClient } from "@angular/common/http";

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
    "galleryNo",
    "name",
    "userNo",
    "createTime",
    "updateTime",
    "createBy",
  ];
  readonly MOBILE_COLUMNS = ["galleryNo", "name", "userNo"];

  @ViewChild("paginator", { static: true })
  paginator: MatPaginator;

  pagingController: PagingController;
  galleries: Gallery[] = [];
  isMobile: boolean = isMobile();
  expandedElement: Gallery = null;
  newGalleryName: string = "";
  showCreateGalleryDiv: boolean = false;

  idEquals = isIdEqual;
  getExpandedEle = (row) => getExpanded(row, this.expandedElement);

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    private navigation: NavigationService,
    private dialog: MatDialog
  ) {}

  ngOnInit() {}

  fetchGalleries() {
    this.http
      .post<any>(`${environment.vfm}/open/api/gallery/list`, {
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          this.pagingController.onTotalChanged(resp.data.paging);
          this.galleries = resp.data.payload;
          this.expandedElement = null;
        },
      });
  }

  createGallery() {
    if (!this.newGalleryName) {
      this.toaster.toast("Please enter new gallery's name");
      return;
    }

    this.http
      .post<any>(`${environment.vfm}/open/api/gallery/new`, {
        name: this.newGalleryName,
      })
      .subscribe({
        next: (resp) => {
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
          .post<any>(`${environment.vfm}/open/api/gallery/delete`, {
            galleryNo: galleryNo,
          })
          .subscribe({
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

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchGalleries();
    this.fetchGalleries();
  }

  updateGallery(galleryNo: string, name: string) {
    if (!galleryNo || !name) return;

    this.http
      .post(`${environment.vfm}/open/api/gallery/update`, {
        galleryNo: galleryNo,
        name: name,
      })
      .subscribe({
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
}
