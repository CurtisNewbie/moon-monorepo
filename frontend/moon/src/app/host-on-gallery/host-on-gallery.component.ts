import { Component, Inject, OnInit } from "@angular/core";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { ConfirmDialog } from "src/common/dialog";
import { filterAlike } from "src/common/select-util";
import { GalleryBrief } from "src/common/gallery";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";

type GlFile = {
  fileKey: string;
  name: string;
  type: string;
};

type Data = {
  files: GlFile[];
};

@Component({
  selector: "app-host-on-gallery",
  templateUrl: "./host-on-gallery.component.html",
  styleUrls: ["./host-on-gallery.component.css"],
})
export class HostOnGalleryComponent implements OnInit {
  /** list of brief info of all galleries that we created */
  galleryBriefs: GalleryBrief[] = [];
  /** name of gallery that we may transfer files to */
  addToGalleryName: string = null;
  /** Auto complete for gallery that we may transfer files to */
  autoCompAddToGalleryName: string[];

  onAddToGalleryNameChanged = () =>
    (this.autoCompAddToGalleryName = filterAlike(
      this.galleryBriefs.map((v) => v.name),
      this.addToGalleryName
    ));

  constructor(
    public dialogRef: MatDialogRef<HostOnGalleryComponent, Data>,
    @Inject(MAT_DIALOG_DATA) public dat: Data,
    private http: HttpClient,
    private confirmDialog: ConfirmDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this._fetchOwnedGalleryBrief();
  }

  transferSelectedToGallery() {
    const addToGalleryNo = this.extractToGalleryNo();
    if (!addToGalleryNo) return;

    let params = this.dat.files.map((f) => {
      return {
        fileKey: f.fileKey,
        galleryNo: addToGalleryNo,
      };
    });

    this.http
      .post<any>(`vfm/open/api/gallery/image/transfer`, {
        images: params,
      })
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.snackBar.open("Request success! It may take a while.", "ok", {
            duration: 3000,
          });
        },
      });
  }

  private extractToGalleryNo(): string {
    const gname = this.addToGalleryName;
    if (!gname) {
      this.snackBar.open("Please select gallery", "ok", { duration: 3000 });
      return;
    }

    let matched: GalleryBrief[] = this.galleryBriefs.filter(
      (v) => v.name === gname
    );
    if (!matched || matched.length < 1) {
      this.snackBar.open(
        "Gallery not found, please check and try again",
        "ok",
        { duration: 3000 }
      );
      return null;
    }
    if (matched.length > 1) {
      this.snackBar.open(
        "Found multiple galleries with the same name, please try again",
        "ok",
        { duration: 3000 }
      );
      return null;
    }
    return matched[0].galleryNo;
  }

  private _fetchOwnedGalleryBrief() {
    this.http.get<any>(`vfm/open/api/gallery/brief/owned`).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.galleryBriefs = resp.data;
        this.onAddToGalleryNameChanged();
      },
    });
  }
}
