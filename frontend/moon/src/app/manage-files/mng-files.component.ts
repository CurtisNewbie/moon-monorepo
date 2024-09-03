import { HttpClient, HttpEventType } from "@angular/common/http";
import {
  Component,
  DoCheck,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from "@angular/core";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";

import {
  emptyUploadFileParam,
  FileInfo,
  FileType,
  SearchFileInfoParam,
  UploadFileParam,
  getFileTypeOpts,
} from "src/common/file-info";
import { PagingController } from "src/common/paging";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";
import { Toaster } from "../notification.service";
import { UserService } from "../user.service";
import { animateElementExpanding, isIdEqual } from "../../animate/animate-util";
import { FileInfoService, TokenType } from "../file-info.service";
import { NavigationService } from "../navigation.service";
import { isMobile } from "src/common/env-util";
import { environment } from "src/environments/environment";
import { ActivatedRoute } from "@angular/router";
import { ImageViewerComponent } from "../image-viewer/image-viewer.component";
import {
  isImage,
  isImageByName,
  isPdf,
  isStreamableVideo,
  guessFileThumbnail,
  isTxt,
  resolveSize,
  isWebpage,
} from "src/common/file";
import { MediaStreamerComponent } from "../media-streamer/media-streamer.component";
import { Option } from "src/common/select-util";
import { isEnterKey } from "src/common/condition";
import { NavType } from "../routes";
import { VfolderAddFileComponent } from "../vfolder-add-file/vfolder-add-file.component";
import { HostOnGalleryComponent } from "../host-on-gallery/host-on-gallery.component";
import { DirectoryMoveFileComponent } from "../directory-move-file/directory-move-file.component";
import { ShareFileQrcodeDialogComponent } from "../share-file-qrcode-dialog/share-file-qrcode-dialog.component";
import { Subscription } from "rxjs";

@Component({
  selector: "app-mng-files",
  templateUrl: "./mng-files.component.html",
  styleUrls: ["./mng-files.component.css"],
  animations: [animateElementExpanding()],
})
export class MngFilesComponent implements OnInit, OnDestroy, DoCheck {
  readonly desktopColumns = [
    "selected",
    "thumbnail",
    "name",
    "parentFileName",
    "uploadTime",
    "size",
    "operation",
  ];
  readonly desktopFolderColumns = [
    "thumbnail",
    "name",
    "uploader",
    "uploadTime",
    "size",
    "operation",
  ];
  readonly mobileColumns = ["fileType", "thumbnail", "name", "operation"];

  allFileTypeOpts: Option<FileType>[] = [];
  guessFileThumbnail = guessFileThumbnail;

  /** expanded fileInfo */
  curr: FileInfo;
  /** expanded fileInfo's id or -1 */
  currId: number = -1;

  /** list of files fetched */
  fileInfoList: FileInfo[] = [];
  /** searching param */
  searchParam: SearchFileInfoParam = {};
  /** controller for pagination */
  pagingController: PagingController;
  /** progress string */
  progress: string = null;
  /** whether current user is using mobile device */
  isMobile: boolean = false;
  /** check if all files are selected */
  isAllSelected: boolean = false;
  /** selected file count */
  selectedCount: number = 0;
  /** is any file selected */
  anySelected: boolean = false;
  /** currently displayed columns */
  displayedColumns: string[] = this._selectColumns();

  // isImage = (f: FileInfo): boolean => this._isImage(f);
  idEquals = isIdEqual;

  // getExpandedEle = (row): FileInfo => getExpanded(row, this.curr, this.isMobile);
  selectExpanded = (row): FileInfo => {
    if (this.isMobile) return null;
    // null means row is the expanded one, so we return null to make it collapsed
    this.curr = this.currId > -1 && row.id == this.currId ? null : { ...row };
    this.currId = this.curr ? this.curr.id : -1;
  };

  isEnterKeyPressed = isEnterKey;
  inSensitiveMode = false;

  /*
  -----------------------

  Virtual Folders

  -----------------------
  */

  /** the folderNo of the folder that we are currently in */
  inFolderNo: string = "";
  /** the name of the folder that we are currently in */
  inFolderName: string = "";

  /*
  -----------------------

  Directory

  -----------------------
  */

  /** the name of the directory that we are currently in */
  inDirFileName: string = null;
  /** the file key of the directory that we are currently in */
  inDirFileKey: string = null;

  /** whether we are making directory */
  makingDir: boolean = false;
  /** name of new dir */
  newDirName: string = null;

  /*
  -----------------------

  Uploading

  -----------------------
  */
  /** whther the upload panel is expanded */
  expandUploadPanel = false;
  /** params for uploading */
  uploadParam: UploadFileParam = emptyUploadFileParam();
  /** displayed upload file name */
  displayedUploadName: string = null;
  /** whether we are uploading */
  isUploading: boolean = false;
  /** name of directory that we may upload files into */
  uploadDirName: string = null;
  /** auto complete for dirs that we may upload file into */
  autoCompUploadDirs: string[] = [];
  /** Always points to current file, so the next will be uploadIndex+1 */
  uploadIndex = -1;
  /** subscription of current uploading */
  uploadSub: Subscription = null;
  /** Ignore upload on duplicate name found*/
  ignoreOnDupName: boolean = true;

  /*
  ----------------------------------

  Labels

  ----------------------------------
  */
  refreshLabel = () => {
    this.allFileTypeOpts = getFileTypeOpts(true);
  };

  @ViewChild("uploadFileInput")
  uploadFileInput: ElementRef;

  setSearchFileType = (fileType) => (this.searchParam.fileType = fileType);

  constructor(
    private userService: UserService,
    private toaster: Toaster,
    private dialog: MatDialog,
    private fileService: FileInfoService,
    private nav: NavigationService,
    private http: HttpClient,
    private route: ActivatedRoute
  ) {}

  ngDoCheck(): void {
    this.anySelected = this.selectedCount > 0;
    this.displayedColumns = this._selectColumns();
  }

  ngOnDestroy(): void {}

  ngOnInit() {
    this.refreshLabel();
    this.isMobile = isMobile();

    this.route.paramMap.subscribe((params) => {
      // vfolder
      this.inFolderNo = params.get("folderNo");
      this.inFolderName = params.get("folderName");

      // directory
      this.inDirFileName = params.get("parentDirName");
      this.inDirFileKey = params.get("parentDirKey");

      // if we are already in a directory, by default we upload to current directory
      if (this.expandUploadPanel && this.inDirFileName) {
        this.uploadDirName = this.inDirFileName;
      }

      if (this.pagingController) {
        if (!this.pagingController.atFirstPage()) {
          this.pagingController.firstPage(); // this also triggers fetchFileInfoList
          // console.log("ngOnInit.firstPage", time())
        } else {
          this.fetchFileInfoList();
          // console.log("ngOnInit.fetchFileInfoList", time())
        }
      }
    });
  }

  // make dir
  mkdir() {
    const dirName = this.newDirName;
    if (!dirName) {
      this.toaster.toast("Please enter new directory name");
      return;
    }

    this.newDirName = null;
    this.http
      .post(`${environment.vfm}/open/api/file/make-dir`, {
        name: dirName,
        parentFile: this.inDirFileKey,
      })
      .subscribe({
        next: () => {
          this.fetchFileInfoList();
          this.makingDir = false;
        },
      });
  }

  // Go to dir, i.e., list files under the directory
  goToDir(name, fileKey) {
    this.expandUploadPanel = false;
    this.curr = null;
    this.resetSearchParam(false, false);
    this.nav.navigateTo(NavType.MANAGE_FILES, [
      { parentDirName: name, parentDirKey: fileKey },
    ]);
  }

  // TODO
  // Move selected to dir
  moveSelectedToDir(into: boolean = true) {
    const selected = this.filterSelected();
    if (!selected || selected.length < 1) {
      this.toaster.toast("Please select files first");
      return;
    }

    if (!into) {
      let msgs = [
        "You sure you want to move these files out of current directory?",
        "",
      ];
      let c = 0;
      for (let f of selected) {
        msgs.push(` ${++c}. ${f.name}`);
      }

      this.dialog
        .open(ConfirmDialogComponent, {
          width: "500px",
          data: {
            title: "Move Files",
            msg: msgs,
            isNoBtnDisplayed: true,
          },
        })
        .afterClosed()
        .subscribe((confirm) => {
          if (!confirm) return;
          this._moveEachToDir(selected, "", 0);
        });
      return;
    }

    this.dialog
      .open(DirectoryMoveFileComponent, {
        width: "500px",
        data: {
          files: selected.map((f, i) => {
            return { name: `${i + 1}. ${f.name}`, fileKey: f.uuid };
          }),
        },
      })
      .afterClosed()
      .subscribe(() => this.fetchFileInfoList());
  }

  private _moveEachToDir(
    selected: FileInfo[],
    dirFileKey: string,
    offset: number
  ) {
    if (offset >= selected.length) {
      this.fetchFileInfoList();
      return;
    }

    let curr = selected[offset];
    this.http
      .post(`${environment.vfm}/open/api/file/move-to-dir`, {
        uuid: curr.uuid,
        parentFileUuid: dirFileKey,
      })
      .subscribe({
        next: (resp) => {
          this._moveEachToDir(selected, dirFileKey, offset + 1);
        },
      });
  }

  /** fetch file info list */
  fetchFileInfoList() {
    this.searchParam.parentFile = this.inDirFileKey;

    this.http
      .post<any>(`${environment.vfm}/open/api/file/list`, {
        paging: this.pagingController.paging,
        filename: this.searchParam.name,
        folderNo: this.inFolderNo,
        parentFile: this.searchParam.parentFile,
        fileType: this.searchParam.fileType,
        sensitive: this.inSensitiveMode,
      })
      .subscribe({
        next: (resp) => {
          this.fileInfoList = [];
          if (resp.data.payload) {
            for (let f of resp.data.payload) {
              f.isFile = f.fileType == FileType.FILE;
              f.isDir = !f.isFile;
              f.fileTypeLabel = f.isFile ? "File" : "Directory";
              f.sizeLabel = resolveSize(f.sizeInBytes);
              f.isDisplayable = this.isDisplayable(f);
              if (f.updateTime) f.updateTime = new Date(f.updateTime);
              if (f.uploadTime) f.uploadTime = new Date(f.uploadTime);
              this.fileInfoList.push(f);

              if (f.thumbnailToken) {
                f.thumbnailUrl =
                  environment.fstore +
                  "/file/raw?key=" +
                  encodeURIComponent(f.thumbnailToken);
              }
            }
          }

          this.pagingController.onTotalChanged(resp.data.paging);
          this.isAllSelected = false;
          this.selectedCount = 0;
        },
        error: (err) => console.log(err),
      });
  }

  /** Upload file */
  upload(): void {
    if (this.isUploading) {
      this.toaster.toast("Uploading, please wait for a moment");
      return;
    }

    if (this.uploadParam.files.length < 1) {
      this.toaster.toast("Please select a file to upload");
      return;
    }

    let isSingleUpload = this._isSingleUpload();

    // single file upload name is required
    if (!this.displayedUploadName && isSingleUpload) {
      this.toaster.toast("Please enter filename");
      return;
    }

    this.uploadParam.ignoreOnDupName = this.ignoreOnDupName;

    if (isSingleUpload) {
      // only need to upload a single file
      this.isUploading = true;
      this.uploadParam.fileName = this.displayedUploadName;
      this._doUpload(this.uploadParam);
    } else {
      // upload one by one
      this.isUploading = true;
      this._doUpload(this._prepNextUpload(), false);
    }
  }

  leaveFolder() {
    if (!this.inFolderNo) return;

    this.nav.navigateTo(NavType.FOLDERS);
  }

  /** Handle events on file selected/changed */
  onFileSelected(files: File[]): void {
    if (this.isUploading) return; // files can't be changed while uploading

    if (files.length < 1) {
      this._resetFileUploadParam();
      return;
    }

    this.uploadParam.files = files;
    this._setDisplayedFileName();

    if (!environment.production) {
      console.log("uploadParam.files", this.uploadParam.files);
    }
  }

  goPrevDir() {
    if (!this.inDirFileKey || !this.inDirFileName) {
      this.inDirFileKey = null;
      this.inDirFileName = null;
      return;
    }

    this.expandUploadPanel = false;
    this.http
      .get<any>(
        `${environment.vfm}/open/api/file/parent?fileKey=${this.inDirFileKey}`
      )
      .subscribe({
        next: (resp) => {
          // console.log("fetchParentFileKey", resp)
          if (resp.data) {
            this.goToDir(resp.data.fileName, resp.data.fileKey);
          } else {
            this.nav.navigateTo(NavType.MANAGE_FILES, []);
          }
        },
      });
  }

  /** Reset all parameters used for searching, and the fetch the list */
  resetSearchParam(
    setFirstPage: boolean = true,
    fetchFileInfoList: boolean = true
  ): void {
    this.curr = null;
    this.currId = -1;

    this.searchParam = {};
    if (setFirstPage && !this.pagingController.atFirstPage()) {
      this.pagingController.firstPage(); // this also triggers fetchFileInfoList
      // console.log("resetSearchParam.firstPage", time())
    } else {
      if (fetchFileInfoList) this.fetchFileInfoList();
    }
  }

  truncateDir(f: FileInfo): void {
    if (!f) {
      return;
    }

    let msgs = [
      `You sure you want to truncate directory '${f.name}'?`,
      "All files in this directory will be deleted.",
    ];
    this.dialog
      .open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Truncate Directory",
          msg: msgs,
          isNoBtnDisplayed: true,
        },
      })
      .afterClosed()
      .subscribe((confirm) => {
        if (!confirm) {
          return;
        }

        this.http
          .post<any>(`${environment.vfm}/open/api/file/dir/truncate`, {
            uuid: f.uuid,
          })
          .subscribe((resp) => {
            this.toaster.toast("Truncating directory, please wait for a while");
            this.fetchFileInfoList();
          });
      });
  }

  deleteSelected(): void {
    let selected = this.fileInfoList
      .filter((f) => f._selected)
      .map((f, i) => {
        return { name: `${i + 1}. ${f.name}`, fileKey: f.uuid };
      });

    if (!selected || selected.length < 1) {
      this.toaster.toast("Select files first");
      return;
    }

    let msgs = [`You sure you want to delete the following files?`, ""];
    for (let s of selected) {
      msgs.push(s.name);
    }

    this.dialog
      .open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Delete Files",
          msg: msgs,
          isNoBtnDisplayed: true,
        },
      })
      .afterClosed()
      .subscribe((confirm) => {
        console.log(confirm);
        if (!confirm) {
          return;
        }
        let fks = [];
        for (let f of selected) {
          fks.push(f.fileKey);
        }

        this.http
          .post<any>(`${environment.vfm}/open/api/file/delete/batch`, {
            fileKeys: fks,
          })
          .subscribe((resp) => {
            this.fetchFileInfoList();
            console.log("deleted", fks);
          });
      });
  }

  /**
   * Delete file
   */
  deleteFile(uuid: string, name: string): void {
    const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
      this.dialog.open(ConfirmDialogComponent, {
        width: "500px",
        data: {
          title: "Delete File",
          msg: [`You sure you want to delete '${name}'`],
          isNoBtnDisplayed: true,
        },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      // console.log(confirm);
      if (confirm) {
        this.http
          .post<any>(`${environment.vfm}/open/api/file/delete`, { uuid: uuid })
          .subscribe((resp) => {
            this.fetchFileInfoList();
          });
      }
    });
  }

  subSetToStr(set: Set<string>, maxCount: number): string {
    let s: string = "";
    let i: number = 0;
    for (let e of set) {
      if (i++ >= maxCount) break;

      s += e + ", ";
    }
    return s.substring(0, s.length - ", ".length);
  }

  /** Cancel the file uploading */
  cancelFileUpload(): void {
    if (!this.isUploading) return;

    if (this.uploadSub != null && !this.uploadSub.closed) {
      this.uploadSub.unsubscribe();
      return;
    }

    this.isUploading = false;
    this._resetFileUploadParam();
    this.toaster.toast("File uploading cancelled");
  }

  /** Update file's info */
  update(u: FileInfo): void {
    if (!u) return;

    this.http
      .post<any>(`${environment.vfm}/open/api/file/info/update`, {
        id: u.id,
        name: u.name,
        sensitiveMode: u.sensitiveMode,
      })
      .subscribe({
        complete: () => {
          this.fetchFileInfoList();
          this.curr = null;
          this.currId = 0;
        },
      });
  }

  /** Guess whether the file is displayable by its name */
  isDisplayable(f: FileInfo): boolean {
    if (!f || !f.isFile) return false;

    const filename: string = f.name;
    if (!filename) return false;

    return (
      isPdf(filename) ||
      isImageByName(filename) ||
      isStreamableVideo(filename) ||
      isTxt(filename) ||
      isWebpage(filename)
    );
  }

  /** Display the file */
  preview(u: FileInfo): void {
    const isStreaming = isStreamableVideo(u.name);
    this.fileService
      .generateFileTempToken(
        u.uuid,
        isStreaming ? TokenType.STREAMING : TokenType.DOWNLOAD
      )
      .subscribe({
        next: (resp) => {
          const token = resp.data;

          const getDownloadUrl = () =>
            environment.fstore + "/file/raw?key=" + encodeURIComponent(token);
          const getStreamingUrl = () =>
            environment.fstore +
            "/file/stream?key=" +
            encodeURIComponent(token);

          if (isStreaming) {
            this.dialog.open(MediaStreamerComponent, {
              data: {
                name: u.name,
                url: getStreamingUrl(),
                token: token,
              },
            });
          } else if (isPdf(u.name)) {
            this.nav.navigateTo(NavType.PDF_VIEWER, [
              { name: u.name, url: getDownloadUrl(), uuid: u.uuid },
            ]);
          } else if (isTxt(u.name)) {
            this.nav.navigateTo(NavType.TXT_VIEWER, [
              { name: u.name, url: getDownloadUrl(), uuid: u.uuid },
            ]);
          } else if (isWebpage(u.name)) {
            this.nav.navigateTo(NavType.WEBPAGE_VIEWER, [
              { name: u.name, url: getDownloadUrl(), uuid: u.uuid },
            ]);
          } else {
            // image
            this.dialog.open(ImageViewerComponent, {
              data: {
                name: u.name,
                url: getDownloadUrl(),
                isMobile: this.isMobile,
                rotate: false,
              },
            });
          }
        },
      });
  }

  generateTempTokenQrCode(fi: FileInfo): void {
    if (!fi) return;

    this.fileService.generateFileTempToken(fi.uuid).subscribe({
      next: (resp) => {
        const dialogRef: MatDialogRef<ShareFileQrcodeDialogComponent, boolean> =
          this.dialog.open(ShareFileQrcodeDialogComponent, {
            data: {
              title: "Share File By QRCode",
              msg: ["Scan QRCode to download the file"],
              img:
                window.location.protocol +
                "//" +
                window.location.host +
                "/" +
                environment.vfm +
                "/open/api/file/token/qrcode?token=" +
                encodeURIComponent(resp.data),
            },
          });

        dialogRef.afterClosed().subscribe((confirm) => {
          // do nothing
        });
      },
    });
  }

  /**
   * Fetch download url and open it in a new tab
   */
  jumpToDownloadUrl(fileKey: string): void {
    this.fileService.jumpToDownloadUrl(fileKey);
  }

  isFileNameInputDisabled(): boolean {
    return this.isUploading || this._isMultipleUpload();
  }

  addToVirtualFolder() {
    if (!this.fileInfoList) {
      this.toaster.toast("Please select files first");
      return;
    }

    let selected = this.fileInfoList
      .filter((f) => f._selected)
      .map((f, i) => {
        return { name: `${i + 1}. ${f.name}`, fileKey: f.uuid };
      });

    if (!selected || selected.length < 1) {
      this.toaster.toast("Please select files first");
      return;
    }

    this.dialog
      .open(VfolderAddFileComponent, {
        width: "500px",
        data: { files: selected },
      })
      .afterClosed()
      .subscribe();
  }

  onPagingControllerReady(pagingController: PagingController) {
    this.pagingController = pagingController;
    this.pagingController.onPageChanged = () => this.fetchFileInfoList();
    this.fetchFileInfoList();
  }

  transferSelectedToGallery() {
    let selected = this.filterSelected(
      (f: FileInfo): boolean => isImage(f) || f.isDir
    ).map((f, i) => {
      return { name: `${i + 1}. ${f.name}`, fileKey: f.uuid };
    });

    if (!selected || selected.length < 1) {
      this.toaster.toast("Please select images or directory first");
      return;
    }

    this.dialog
      .open(HostOnGalleryComponent, {
        width: "500px",
        data: { files: selected },
      })
      .afterClosed()
      .subscribe();
  }

  selectFile(event: any, f: FileInfo) {
    const isChecked = event.checked;
    f._selected = isChecked;
    let delta = isChecked ? 1 : -1;
    this.selectedCount += delta;
  }

  selectAll() {
    this.isAllSelected = !this.isAllSelected;
    let total = 0;

    this.fileInfoList.forEach((v) => {
      v._selected = this.isAllSelected;
      total += 1;
    });
    this.selectedCount = this.isAllSelected ? total : 0;
  }

  toggleMkdirPanel() {
    this.makingDir = !this.makingDir;
    if (this.makingDir) {
      this.expandUploadPanel = false;
    }
  }

  toggleUploadPanel() {
    this.expandUploadPanel = !this.expandUploadPanel;

    if (this.expandUploadPanel) {
      this.makingDir = false;

      // if we are already in a directory, by default we upload to current directory
      if (!this.uploadParam.parentFile && this.inDirFileName) {
        this.uploadDirName = this.inDirFileName;
      }
    }
  }

  // -------------------------- private helper methods ------------------------

  private _concatTempFileDownloadUrl(tempToken: string): string {
    return (
      window.location.protocol +
      "//" +
      window.location.host +
      "/" +
      environment.fstore +
      "/file/raw?key=" +
      encodeURIComponent(tempToken)
    );
  }

  private _setDisplayedFileName(): void {
    if (!this.uploadParam || !this.uploadParam.files) return;

    const files = this.uploadParam.files;
    const firstFile: File = files[0];
    if (this._isSingleUpload()) this.displayedUploadName = firstFile.name;
    else
      this.displayedUploadName = `Batch Upload: ${files.length} files in total`;
  }

  private _resetFileUploadParam(): void {
    if (this.isUploading) return;

    this.isAllSelected = false;
    this.uploadParam = emptyUploadFileParam();

    if (this.uploadFileInput) {
      this.uploadFileInput.nativeElement.value = null;
    }

    this.uploadIndex = -1;
    this.displayedUploadName = null;
    this.progress = null;

    if (!this.inDirFileName) {
      this.uploadDirName = null;
    }

    this.pagingController.firstPage();
    this.expandUploadPanel = false;
  }

  private _prepNextUpload(): UploadFileParam {
    if (!this.isUploading) return null;
    if (this._isSingleUpload()) return null;

    let i = this.uploadIndex; // if this is the first one, i will be -1
    let files = this.uploadParam.files;
    let next_i = i + 1;

    if (next_i >= files.length) return null;

    let next = files[next_i];
    if (!next) return null;

    this.uploadIndex = next_i;

    return {
      fileName: next.name,
      files: [next],
      ignoreOnDupName: this.uploadParam.ignoreOnDupName,
    };
  }

  private _updateUploadProgress(
    filename: string,
    loaded: number,
    total: number
  ) {
    // how many files left
    let remaining;
    let index = this.uploadIndex;
    if (index == -1) remaining = "";
    else {
      let files = this.uploadParam.files;
      if (!files) remaining = "";
      else {
        let len = files.length;
        if (index >= len) remaining = "";
        else remaining = `${len - this.uploadIndex - 1} file remaining`;
      }
    }

    // upload progress
    let p = Math.round((100 * loaded) / total).toFixed(2);
    let ps;
    if (p == "100.00") ps = `Processing '${filename}' ... ${remaining} `;
    else ps = `Uploading ${filename} ${p}% ${remaining} `;
    this.progress = ps;
  }

  private _doUpload(
    uploadParam: UploadFileParam,
    fetchOnComplete: boolean = true
  ) {
    uploadParam.parentFile = this.inDirFileKey;
    const onComplete = () => {
      if (fetchOnComplete) setTimeout(() => this.fetchFileInfoList(), 1_000);

      let next = this._prepNextUpload();
      if (!next) {
        this.progress = null;
        this.isUploading = false;
        this._resetFileUploadParam();
        this.fetchFileInfoList();
      } else {
        this._doUpload(next, false); // upload next file
      }
    };

    const abortUpload = () => {
      this.progress = null;
      this.isUploading = false;
      this.toaster.toast(`Failed to upload file ${name} `);
      this._resetFileUploadParam();
    };

    const name = uploadParam.fileName;
    const uploadFileCallback = () => {
      this.uploadSub = this.fileService
        .uploadToMiniFstore(uploadParam)
        .subscribe({
          next: (event) => {
            if (event.type === HttpEventType.UploadProgress) {
              this._updateUploadProgress(
                uploadParam.fileName,
                event.loaded,
                event.total
              );
            }

            // TODO: refactor this later, this is so ugly
            if (event.type == HttpEventType.Response) {
              let fstoreRes = event.body;
              if (fstoreRes.error) {
                abortUpload();
                return;
              }

              // create the record in vfm
              this.http
                .post(`${environment.vfm}/open/api/file/create`, {
                  filename: uploadParam.fileName,
                  fstoreFileId: fstoreRes.data,
                  parentFile: uploadParam.parentFile,
                })
                .subscribe({
                  complete: onComplete,
                  error: () => {
                    abortUpload();
                  },
                });
            }
          },
          error: () => {
            abortUpload();
          },
        });
    };

    if (!uploadParam.ignoreOnDupName) {
      uploadFileCallback();
    } else {
      let pf = uploadParam.parentFile
        ? encodeURIComponent(uploadParam.parentFile)
        : "";

      // preflight check whether the filename exists already
      this.http
        .get<any>(
          `${
            environment.vfm
          }/open/api/file/upload/duplication/preflight?fileName=${encodeURIComponent(
            name
          )}&parentFileKey=${pf}`
        )
        .subscribe({
          next: (resp) => {
            let isDuplicate = resp.data;
            if (!isDuplicate) {
              uploadFileCallback();
            } else {
              this._updateUploadProgress(uploadParam.fileName, 100, 100);

              // skip this file, it exists already
              onComplete();
            }
          },
        });
    }
  }

  private _isSingleUpload() {
    return !this._isMultipleUpload();
  }

  private _isMultipleUpload() {
    return this.uploadParam.files.length > 1;
  }

  private _selectColumns() {
    if (isMobile()) return this.mobileColumns;
    return this.inFolderNo ? this.desktopFolderColumns : this.desktopColumns;
  }

  /** Filter selected files */
  private filterSelected(...predicates): FileInfo[] {
    return this.fileInfoList
      .map((v) => {
        if (!v._selected) return null;
        for (let p of predicates) {
          if (!p(v)) return null;
        }

        return v;
      })
      .filter((v) => v != null);
  }

  onRowClicked(row: FileInfo) {
    if (row.isDir) {
      this.goToDir(row.name, row.uuid);
      return;
    }
    if (row.isDisplayable) {
      this.preview(row);
    }
  }

  sensitiveModeChecked(event, file) {
    file.sensitiveMode = event.checked ? "Y" : "N";
    console.log("checked?", file);
  }

  canUnpack(fi: FileInfo): boolean {
    return fi.name && fi.name.toLowerCase().endsWith(".zip");
  }

  unpack(fi: FileInfo) {
    this.http
      .post(`${environment.vfm}/open/api/file/unpack`, {
        fileKey: fi.uuid,
        parentFileKey: this.inDirFileKey,
      })
      .subscribe({
        next: () => {
          this.fetchFileInfoList();
          this.toaster.toast(`Unpacking ${fi.name}, please be patient.`);
          this.currId = -1;
        },
      });
    return false;
  }

  /**
   * Generate temporary token for downloading
   */
  generateTempToken(u: FileInfo): void {
    if (!u) return;

    this.fileService.generateFileTempToken(u.uuid).subscribe({
      next: (resp) => {
        const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
          this.dialog.open(ConfirmDialogComponent, {
            width: "700px",
            data: {
              title: "Share File",
              msg: [
                "Link to download this file:",
                this._concatTempFileDownloadUrl(resp.data),
              ],
              isNoBtnDisplayed: false,
            },
          });

        dialogRef.afterClosed().subscribe((confirm) => {
          // do nothing
        });
      },
    });
  }
}
