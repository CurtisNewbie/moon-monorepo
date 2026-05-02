import {
  Component,
  EventEmitter,
  Inject,
  OnDestroy,
  OnInit,
  Output,
  ViewEncapsulation,
} from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import {
  Lightbox,
  LightboxConfig,
  LightboxEvent,
  LIGHTBOX_EVENT,
} from "ngx-lightbox";
import { Subscription } from "rxjs";
import { Env } from "src/common/env-util";

export interface ImgViewerDialogData {
  name: string;
  url: string;
  rotate: boolean;
}

@Component({
  selector: "app-image-viewer",
  templateUrl: "./image-viewer.component.html",
  styleUrls: ["./image-viewer.component.css"],
  encapsulation: ViewEncapsulation.None
})
export class ImageViewerComponent implements OnInit, OnDestroy {
  lbsub: Subscription;
  lightboxdiv: any;
  @Output() swipeLeft = new EventEmitter<boolean>();
  @Output() swipeRight = new EventEmitter<boolean>();

  private _destroyed = false;
  private _timeoutIds: number[] = [];
  private _touchStartHandler: ((evt: TouchEvent) => void) | null = null;
  private _touchMoveHandler: ((evt: TouchEvent) => void) | null = null;

  /** setTimeout wrapper that tracks IDs for cleanup on destroy */
  private _setTimeout(fn: () => void, ms: number): number {
    const id = window.setTimeout(() => {
      if (this._destroyed) return;
      // Remove from tracking once fired
      const ix = this._timeoutIds.indexOf(id);
      if (ix > -1) this._timeoutIds.splice(ix, 1);
      fn();
    }, ms);
    this._timeoutIds.push(id);
    return id;
  }

  constructor(
    private env: Env,
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
          if (evt.touches.length > 1) {
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

        // Store handler references for cleanup
        this._touchStartHandler = handleTouchStart;
        this._touchMoveHandler = handleTouchMove;

        this.lightboxdiv.addEventListener(
          "touchstart",
          handleTouchStart,
          false
        );
        this.lightboxdiv.addEventListener("touchmove", handleTouchMove, false);

        let checkCompleteInterval = 50;
        this._setTimeout(() => {
          let imgs = this.lightboxdiv.querySelectorAll("#image");
          console.log(imgs);
          if (imgs && imgs.length > 0) {
            let img = imgs[0];
            const checkComplete = () => {
              if (this._destroyed) return;
              console.log(
                "img",
                img.complete,
                img.clientWidth,
                img.clientHeight
              );
              if (!img.complete) {
                this._setTimeout(checkComplete, checkCompleteInterval);
                return;
              }

              let resizeRetry = 0;
              const resizeRetryInterval = 50;
              const resizeMaxRetry = 10;
              const resize = () => {
                if (this._destroyed) return;
                console.log(
                  "img check zoom-in",
                  img.complete,
                  " width ",
                  img.clientWidth,
                  " height ",
                  img.clientHeight
                );
                resizeRetry += 1;

                // image rendered
                if (img.clientWidth > 0) {
                  const viewportWidth = window.innerWidth;
                  const currentWidth = img.clientWidth;
                  const currentHeight = img.clientHeight;
                  let targetWidth = currentWidth;
                  let scale = 1.0;

                  // Check if it's a long portrait image
                  if (this.isLongPortraitImage(img)) {
                    // Long portrait images
                    if (viewportWidth < 600) {
                      // Mobile: viewport-aware zoom (85-90%)
                      const targetPercentage = viewportWidth < 400 ? 0.85 : 0.90;
                      targetWidth = viewportWidth * targetPercentage;
                      scale = targetWidth / currentWidth;
                    } else {
                      // Desktop: 40-50% viewport width
                      const targetPercentage = 0.40;
                      targetWidth = viewportWidth * targetPercentage;
                      scale = targetWidth / currentWidth;
                    }
                  } else {
                    // Normal images
                    if (viewportWidth < 600) {
                      // Mobile: 60-70% viewport width (constrained)
                      const targetPercentage = viewportWidth < 400 ? 0.60 : 0.70;
                      targetWidth = viewportWidth * targetPercentage;
                      scale = targetWidth / currentWidth;
                    } else {
                      // Desktop: keep pixel thresholds (only scale truly tiny images)
                      if (currentWidth < 150) {
                        scale = 2.5;
                      }
                      targetWidth = currentWidth * scale;
                    }
                  }

                  // Apply width adjustment if scaling is needed
                  if (scale > 1.0) {
                    const newHeight = currentHeight * scale;
                    img.style.width = targetWidth + 'px';
                    img.style.height = newHeight + 'px';
                    img.style.maxWidth = 'none';

                    // Update container width to match scaled image
                    const container = this.lightboxdiv.querySelector('.lb-outerContainer');
                    if (container) {
                      container.style.width = targetWidth + 'px';
                    }
                  }
                } else if (resizeRetry < resizeMaxRetry) {
                  // image not rendered, try again later
                  this._setTimeout(resize, resizeRetryInterval);
                }
              };

              this._setTimeout(resize, resizeRetryInterval);
            };

            checkComplete();
          }
        }, checkCompleteInterval);
      }
      return ele;
    };
    _lbConfig.wrapAround = false;
    _lbConfig.disableScrolling = true;
    _lbConfig.showZoom = true;
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

  ngOnDestroy(): void {
    this._destroyed = true;

    // Clear all pending timeouts
    for (const id of this._timeoutIds) {
      window.clearTimeout(id);
    }
    this._timeoutIds = [];

    // Remove touch event listeners
    if (this.lightboxdiv) {
      if (this._touchStartHandler) {
        this.lightboxdiv.removeEventListener("touchstart", this._touchStartHandler, false);
        this._touchStartHandler = null;
      }
      if (this._touchMoveHandler) {
        this.lightboxdiv.removeEventListener("touchmove", this._touchMoveHandler, false);
        this._touchMoveHandler = null;
      }
    }

    if (this.lbsub && !this.lbsub.closed) {
      this.lbsub.unsubscribe();
    }
  }

  ngOnInit() {}

  getTouches(evt) {
    return (
      evt.touches || // browser API
      evt.originalEvent.touches
    ); // jQuery
  }

  private isLongPortraitImage(img: HTMLImageElement): boolean {
    const width = img.naturalWidth;
    const height = img.naturalHeight;
    const ratio = height / width;

    // Check if natural dimensions are available
    if (!width || !height || width === 0 || height === 0) {
      return false;
    }

    // Long portrait criteria:
    // 1. Height/width ratio >= 3
    // OR
    // 2. Height > 2000px AND ratio > 2
    return ratio >= 3 || (height > 2000 && ratio > 2);
  }
}
