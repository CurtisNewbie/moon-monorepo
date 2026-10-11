import {
  Component,
  ElementRef,
  EventEmitter,
  Inject,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";

export interface MediaStreamerDialogData {
  name: string;
  url: string;
  token: string;
}

@Component({
  selector: "app-media-streamer",
  templateUrl: "./media-streamer.component.html",
  styleUrls: ["./media-streamer.component.css"],
})
export class MediaStreamerComponent implements OnInit, OnDestroy {
  name: string;
  token: string;
  srcUrl: string;
  loadFailed = false;

  @ViewChild("videoPlayer", { static: false })
  videoPlayer?: ElementRef<HTMLVideoElement>;

  @Output("prev")
  prevEventEmitter = new EventEmitter();

  @Output("next")
  nextEventEmitter = new EventEmitter();

  constructor(
    public dialogRef: MatDialogRef<
      MediaStreamerComponent,
      MediaStreamerDialogData
    >,
    @Inject(MAT_DIALOG_DATA) public data: MediaStreamerDialogData
  ) {}

  ngOnDestroy(): void {
    const v = this.videoPlayer?.nativeElement as HTMLVideoElement | undefined;
    if (v) {
      v.pause();
      v.removeAttribute("src");
      v.load();
    }
  }

  ngOnInit() {
    this.name = this.data.name;
    this.srcUrl = location.origin + "/" + this.data.url;
    this.token = this.data.token;
  }

  onVideoError(): void {
    this.loadFailed = true;
  }

  prev() {
    this.prevEventEmitter.emit();
  }

  next() {
    this.nextEventEmitter.emit();
  }
}
