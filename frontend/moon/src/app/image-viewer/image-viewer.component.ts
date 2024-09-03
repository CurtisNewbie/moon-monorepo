import { Component, Inject, OnDestroy, OnInit } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { Lightbox, LightboxConfig, LightboxEvent, LIGHTBOX_EVENT } from "ngx-lightbox";
import { Subscription } from "rxjs";

export interface ImgViewerDialogData {
  name: string;
  url: string;
  isMobile: boolean;
  rotate: boolean;
}

@Component({
  selector: "app-image-viewer",
  templateUrl: "./image-viewer.component.html",
  styleUrls: ["./image-viewer.component.css"],
})
export class ImageViewerComponent implements OnInit, OnDestroy {
  lbsub: Subscription;

  constructor(
    private _lightbox: Lightbox,
    private _lbConfig: LightboxConfig,
    private _lightboxEvent: LightboxEvent,
    public dialogRef: MatDialogRef<ImageViewerComponent, ImgViewerDialogData>,
    @Inject(MAT_DIALOG_DATA) public data: ImgViewerDialogData
  ) {
    _lbConfig.containerElementResolver = (doc: Document) =>
      doc.getElementById("lightboxdiv");
    _lbConfig.wrapAround = false;
    _lbConfig.disableScrolling = true;
    _lbConfig.showZoom = false;
    _lbConfig.resizeDuration = 0.1;
    _lbConfig.fadeDuration = 0.1;
    _lbConfig.showRotate = data.rotate;
    _lbConfig.fitImageInViewPort = true;
    _lbConfig.showImageNumberLabel = false;
    _lbConfig.centerVertically = true;

    this.lbsub = this._lightboxEvent.lightboxEvent$.subscribe((evt: any) => {
      if (evt.id === LIGHTBOX_EVENT.CLOSE) {
        this.lbsub.unsubscribe();
        this.dialogRef.close();
      }
    })

    this._lightbox.open([{
      src: this.data.url,
      thumb: this.data.url,
      downloadUrl: this.data.url
    }], 0, {});
  }

  ngOnDestroy(): void {
  }

  ngOnInit() {
  }

}
