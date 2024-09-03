import { Component, OnInit, ViewChild, ViewEncapsulation } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { HttpClient } from "@angular/common/http";
import { PagingController } from "src/common/paging";
import { ListGalleryImagesResp } from "src/common/gallery";
import { environment } from "src/environments/environment";
import { Resp } from "src/common/resp";
import { NavigationService } from "../navigation.service";
import { IAlbum, Lightbox, LightboxConfig } from "ngx-lightbox";
import { NavType } from "../routes";

@Component({
  selector: "app-gallery-image",
  templateUrl: "./gallery-image.component.html",
  styleUrls: ["./gallery-image.component.css"],
  encapsulation: ViewEncapsulation.None,
})
export class GalleryImageComponent implements OnInit {
  pagingController: PagingController;
  galleryNo: string = null;
  title = "fantahsea";
  images: IAlbum[] = [];

  constructor(
    private route: ActivatedRoute,
    private http: HttpClient,
    private navigation: NavigationService,
    private _lightbox: Lightbox,
    private _lbConfig: LightboxConfig
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
  }

  fetchImages(): void {
    if (!this.galleryNo) this.navigation.navigateTo(NavType.GALLERY);

    this.http
      .post<Resp<ListGalleryImagesResp>>(
        `${environment.vfm}/open/api/gallery/images`,
        { galleryNo: this.galleryNo, paging: this.pagingController.paging }
      )
      .subscribe({
        next: (resp) => {
          this.pagingController.onTotalChanged(resp.data.paging);

          this.images = [];
          if (resp.data.images) {
            let imgs = resp.data.images;
            this.images = [];
            for (let i = 0; i < imgs.length; i++) {
              let src =
                environment.fstore +
                "/file/raw?key=" +
                encodeURIComponent(imgs[i].fileTempToken);
              let thumb =
                environment.fstore +
                "/file/raw?key=" +
                encodeURIComponent(imgs[i].thumbnailToken);
              this.images.push({
                src: src,
                thumb: thumb,
                downloadUrl: src,
              });
            }
          }
        },
      });
  }

  open(index: number): void {
    this._lightbox.open(this.images, index, {
      wrapAround: true,
      showImageNumberLabel: true,
    });
  }

  close(): void {
    // close lightbox programmatically
    this._lightbox.close();
  }

  onPagingControllerReady(pc: any) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchImages();
    this.pagingController.setPageLimit(40);
    this.pagingController.PAGE_LIMIT_OPTIONS = [20, 40, 60, 100, 500];

    this.fetchImages();
  }
}
