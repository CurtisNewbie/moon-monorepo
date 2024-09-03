import { HttpClient } from "@angular/common/http";
import {
  Component,
  DoCheck,
  OnDestroy,
  OnInit,
  ViewChild,
} from "@angular/core";
import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import {
  MatSelectionList,
  MatSelectionListChange,
} from "@angular/material/list";
import { Subscription } from "rxjs";
import { environment } from "src/environments/environment";
import { VFolder } from "src/common/folder";
import { PagingController } from "src/common/paging";
import { Resp } from "src/common/resp";
import { UserInfo } from "src/common/user-info";
import { GrantAccessDialogComponent } from "../grant-access-dialog/grant-access-dialog.component";
import { NavigationService } from "../navigation.service";
import { Toaster } from "../notification.service";
import { NavType } from "../routes";
import { UserService } from "../user.service";
import { isEnterKey } from "src/common/condition";
import { ConfirmDialogComponent } from "../dialog/confirm/confirm-dialog.component";

@Component({
  selector: "app-folder",
  templateUrl: "./folder.component.html",
  styleUrls: ["./folder.component.css"],
})
export class FolderComponent implements OnInit, DoCheck, OnDestroy {
  user: UserInfo;
  userSub: Subscription;
  pagingController: PagingController;
  newFolderName: string = "";
  creatingFolder: boolean = false;
  searchParam = {
    name: "",
    paging: null,
  };
  folders: VFolder[] = [];
  selected: VFolder[] = [];
  onEnterPressed = isEnterKey;
  isOneSelected: boolean;

  @ViewChild("folderList")
  folderList: MatSelectionList;

  constructor(
    private http: HttpClient,
    private toaster: Toaster,
    private navi: NavigationService,
    private dialog: MatDialog,
    private userService: UserService
  ) {}

  ngOnInit(): void {
    this.userSub = this.userService.userInfoObservable.subscribe((u) => {
      this.user = u;
    });
  }

  ngOnDestroy(): void {
    if (this.userSub) this.userSub.unsubscribe();
  }

  ngDoCheck(): void {
    if (!this.folderList) {
      this.isOneSelected = false;
      return;
    }

    let selected = this.folderList.selectedOptions.selected;
    this.isOneSelected = selected.length == 1;
  }

  isOwner(f: VFolder): boolean {
    return f.createBy == this.user.username;
  }

  popToRemoveVFolder(f: VFolder): void {
    if (!f) return;

    const dialogRef: MatDialogRef<ConfirmDialogComponent, boolean> =
      this.dialog.open(ConfirmDialogComponent, {
        width: "700px",
        data: {
          title: "Delete Virtual Folder",
          msg: [`You sure you want to delete '${f.name}'`],
          isNoBtnDisplayed: true,
        },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      if (!confirm) {
        return;
      }
      this.removeVFolder(f.folderNo);
    });
  }

  removeVFolder(folderNo: string) {
    this.http
      .post<Resp<any>>(`${environment.vfm}/open/api/vfolder/remove`, {
        folderNo: folderNo,
      })
      .subscribe({
        next: (resp) => {
          this.toaster.toast("Virtual Folder Removed");
          this.fetchFolders();
        },
      });
  }

  popToGrantAccess(f: VFolder): void {
    if (!f) return;

    const dialogRef: MatDialogRef<GrantAccessDialogComponent, boolean> =
      this.dialog.open(GrantAccessDialogComponent, {
        width: "700px",
        data: { folderNo: f.folderNo, name: f.name },
      });

    dialogRef.afterClosed().subscribe((confirm) => {
      // do nothing
    });
  }

  selectionChanged(event: MatSelectionListChange): void {
    this.selected = event.options.filter((o) => o.selected).map((o) => o.value);
  }

  selectFolder(f: VFolder): void {
    this.navi.navigateTo(NavType.MANAGE_FILES, [
      { folderNo: f.folderNo, folderName: f.name },
    ]);
  }

  fetchFolders(): void {
    this.searchParam.paging = this.pagingController.paging;
    this.http
      .post<Resp<any>>(
        `${environment.vfm}/open/api/vfolder/list`,
        this.searchParam
      )
      .subscribe({
        next: (resp) => {
          this.folders = [];
          if (resp.data.payload) {
            this.folders = resp.data.payload.map((r) => {
              if (r.createTime) r.createTime = new Date(r.createTime);
              if (r.updateTime) r.updateTime = new Date(r.updateTime);
              return r;
            });
          }
          this.pagingController.onTotalChanged(resp.data.paging);
        },
      });
  }

  resetSearchParam(): void {
    this.searchParam.name = "";
    this.folderList.deselectAll();
    this.fetchFolders();
  }

  createFolder(): void {
    if (!this.newFolderName) {
      this.toaster.toast("Please enter folder name");
      return;
    }

    this.creatingFolder = false;
    this.http
      .post<Resp<void>>(`${environment.vfm}/open/api/vfolder/create`, {
        name: this.newFolderName,
      })
      .subscribe({
        next: (resp) => {
          this.fetchFolders();
          this.newFolderName = "";
        },
      });
  }

  onPagingControllerReady(pc: PagingController) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchFolders();
    this.pagingController.PAGE_LIMIT_OPTIONS = [5, 10];
    this.pagingController.paging.limit =
      this.pagingController.PAGE_LIMIT_OPTIONS[0];
    this.fetchFolders();
  }
}
