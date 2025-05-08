import { Component, OnInit, ViewChild, ViewEncapsulation } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { HttpClient } from "@angular/common/http";
import { ListGalleryImagesResp } from "src/common/gallery";
import { Resp } from "src/common/resp";
import { NavigationService } from "../navigation.service";
import {
  Lightbox,
  LIGHTBOX_EVENT,
  LightboxConfig,
  LightboxEvent,
} from "ngx-lightbox";
import { NavType } from "../routes";
import { MatMenuTrigger } from "@angular/material/menu";
import { BrowseHistoryRecorder } from "src/common/browse-history";
import { Subscription } from "rxjs";
import { MatSnackBar } from "@angular/material/snack-bar";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";

@Component({
  selector: "app-gallery-image",
  templateUrl: "./gallery-image.component.html",
  styleUrls: ["./gallery-image.component.css"],
  encapsulation: ViewEncapsulation.None,
})
export class GalleryImageComponent implements OnInit {
  galleryNo: string = null;
  title = "fantahsea";
  images = [];

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  private lbxSub: Subscription;

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private navigation: NavigationService,
    private _lightbox: Lightbox,
    private _lbConfig: LightboxConfig,
    private _lightboxEvent: LightboxEvent,
    private browseHistoryRecorder: BrowseHistoryRecorder,
    private snackBar: MatSnackBar
  ) {
    _lbConfig.containerElementResolver = (doc: Document) =>
      doc.getElementById("lightboxdiv");
    _lbConfig.wrapAround = false;
    _lbConfig.disableScrolling = false;
    _lbConfig.showZoom = false;
    _lbConfig.resizeDuration = 0.2;
    _lbConfig.fadeDuration = 0.2;
    _lbConfig.showRotate = true;
    _lbConfig.showImageNumberLabel = true;
    _lbConfig.centerVertically = true;
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe((params) => {
      let galleryNo = params.get("galleryNo");
      if (galleryNo) this.galleryNo = galleryNo;
    });
    this.pagingController.setPageLimit(40);
    this.pagingController.PAGE_LIMIT_OPTIONS = [20, 40, 60, 100, 500];
  }

  fetchImages(): void {
    if (!this.galleryNo) this.navigation.navigateTo(NavType.GALLERY);

    this.http
      .post<Resp<ListGalleryImagesResp>>(`vfm/open/api/gallery/images`, {
        galleryNo: this.galleryNo,
        paging: this.pagingController.paging,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.pagingController.onTotalChanged(resp.data.paging);

          this.images = [];
          if (resp.data.images) {
            let imgs = resp.data.images;
            this.images = [];
            for (let i = 0; i < imgs.length; i++) {
              let src =
                "fstore/file/raw?key=" +
                encodeURIComponent(imgs[i].fileTempToken);
              let thumb =
                "fstore/file/raw?key=" +
                encodeURIComponent(imgs[i].thumbnailToken);
              this.images.push({
                src: src,
                thumb: thumb,
                downloadUrl: src,
                fileKey: imgs[i].fileKey,
              });
            }
          }
          console.log("images: ", this.images);
        },
      });
  }

  open(index: number): void {
    this.browseHistoryRecorder.record(this.images[index].fileKey);
    this.lbxSub = this._lightboxEvent.lightboxEvent$.subscribe((event: any) => {
      if (event.id === LIGHTBOX_EVENT.CLOSE) {
        this.lbxSub.unsubscribe();
      }
      if (event.id === LIGHTBOX_EVENT.CHANGE_PAGE) {
        this.browseHistoryRecorder.record(this.images[event.data].fileKey);
      }
    });
    this._lightbox.open(this.images, index, {
      wrapAround: true,
      showImageNumberLabel: true,
    });
  }

  close(): void {
    // close lightbox programmatically
    this._lightbox.close();
  }

  // https://stackoverflow.com/questions/77608499/angular-material-custom-context-menu-right-click
  @ViewChild(MatMenuTrigger, { static: true }) matMenuTrigger: MatMenuTrigger;
  menuTopLeftPosition = { x: "0", y: "0" };
  menuBoundImage = null;

  onRightClick(event, image) {
    event.preventDefault();

    this.menuTopLeftPosition.x = event.clientX + "px";
    this.menuTopLeftPosition.y = event.clientY + "px";
    this.menuBoundImage = image;
    this.matMenuTrigger.openMenu();
  }

  goToFile() {
    if (!this.menuBoundImage) {
      return;
    }
    this.navigation.navigateTo(NavType.MANAGE_FILES, [
      { searchedFileKey: this.menuBoundImage.fileKey },
    ]);
  }
}
