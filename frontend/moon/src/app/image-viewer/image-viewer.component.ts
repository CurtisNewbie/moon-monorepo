import {
  Component,
  EventEmitter,
  Inject,
  OnDestroy,
  OnInit,
  Output,
} from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import {
  Lightbox,
  LightboxConfig,
  LightboxEvent,
  LIGHTBOX_EVENT,
} from "ngx-lightbox";
import { Subscription } from "rxjs";

export interface ImgViewerDialogData {
  name: string;
  url: string;
  rotate: boolean;
}

@Component({
  selector: "app-image-viewer",
  templateUrl: "./image-viewer.component.html",
  styleUrls: ["./image-viewer.component.css"],
})
export class ImageViewerComponent implements OnInit, OnDestroy {
  lbsub: Subscription;
  lightboxdiv: any;
  @Output() swipeLeft = new EventEmitter<boolean>();
  @Output() swipeRight = new EventEmitter<boolean>();

  constructor(
    private _lightbox: Lightbox,
    private _lbConfig: LightboxConfig,
    private _lightboxEvent: LightboxEvent,
    public dialogRef: MatDialogRef<ImageViewerComponent, ImgViewerDialogData>,
    @Inject(MAT_DIALOG_DATA) public data: ImgViewerDialogData
  ) {
    _lbConfig.containerElementResolver = (doc: Document) => {
      let ele = doc.getElementById("lightboxdiv");
      let firstTime = this.lightboxdiv == null;
      this.lightboxdiv = ele;
      if (firstTime) {
        var xDown = null;
        var yDown = null;
        let handleTouchStart = (evt) => {
          const firstTouch = this.getTouches(evt)[0];
          xDown = firstTouch.clientX;
          yDown = firstTouch.clientY;
        };
        let handleTouchMove = (evt) => {
          if (!xDown || !yDown) {
            return;
          }

          var xUp = evt.touches[0].clientX;
          var yUp = evt.touches[0].clientY;

          var xDiff = xDown - xUp;
          var yDiff = yDown - yUp;

          if (Math.abs(xDiff) > Math.abs(yDiff)) {
            if (xDiff > 0) {
              this.swipeRight.emit(true);
            } else {
              this.swipeLeft.emit(true);
            }
          } else {
            if (yDiff > 0) {
              /* down swipe */
            } else {
              /* up swipe */
            }
          }
          /* reset values */
          xDown = null;
          yDown = null;
        };

        this.lightboxdiv.addEventListener(
          "touchstart",
          handleTouchStart,
          false
        );
        this.lightboxdiv.addEventListener("touchmove", handleTouchMove, false);
      }
      return ele;
    };
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
    });

    this._lightbox.open(
      [
        {
          src: this.data.url,
          thumb: this.data.url,
          downloadUrl: this.data.url,
        },
      ],
      0,
      {}
    );
  }

  ngOnDestroy(): void {}

  ngOnInit() {}

  getTouches(evt) {
    return (
      evt.touches || // browser API
      evt.originalEvent.touches
    ); // jQuery
  }
}
