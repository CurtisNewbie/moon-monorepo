import { HttpClient } from "@angular/common/http";
import { Component, OnInit, ViewChild } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { MatSnackBar } from "@angular/material/snack-bar";
import { isEnterKey } from "src/common/condition";
import { Paging } from "src/common/paging";
import { ControlledPaginatorComponent } from "../controlled-paginator/controlled-paginator.component";
import { EditNoteDialogComponent } from "../edit-note-dialog/edit-note-dialog.component";
import { Env } from "src/common/env-util";
import { I18n } from "../i18n.service";

export interface SaveNoteReq {
  title?: string; // Required.
  content?: string; // Required.
}

export interface ListNoteReq {
  keywords?: string;
  paging?: Paging;
}

export interface Note {
  recordId?: string;
  title?: string;
  content?: string;
  userNo?: string;
}

export interface PageRes {
  paging?: Paging;
  payload?: Note[];
}

@Component({
  selector: "app-mng-note",
  templateUrl: "./mng-note.component.html",
  styleUrls: ["./mng-note.component.css"],
})
export class MngNoteComponent implements OnInit {
  addNotePanelDisplayed: boolean = false;
  isEnter = isEnterKey;
  tabdata = [];
  newReq: SaveNoteReq = {};
  listReq: ListNoteReq = {};

  @ViewChild(ControlledPaginatorComponent)
  pagingController: ControlledPaginatorComponent;

  constructor(
    public env: Env,
    private snackBar: MatSnackBar,
    private http: HttpClient,
    private dialog: MatDialog,
    public i18n: I18n
  ) {}

  trl(k) {
    return this.i18n.trl("mng-note", k);
  }

  ngOnInit(): void {
    this.fetchList();
  }

  fetchList() {
    let req: ListNoteReq = this.listReq;
    this.http.post<any>(`/user-vault/open/api/note/list-notes`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        let dat: PageRes = resp.data;
        this.tabdata = dat.payload;
        this.pagingController.onTotalChanged(dat.paging);
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open(this.trl("requestFailedUnknownError"), "ok", {
          duration: 3000,
        });
      },
    });
  }

  saveNote() {
    let req: SaveNoteReq = this.newReq;
    this.http.post<any>(`/user-vault/open/api/note/save-note`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        this.newReq = {};
        this.addNotePanelDisplayed = false;
        this.fetchList();
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open(this.trl("requestFailedUnknownError"), "ok", {
          duration: 3000,
        });
      },
    });
  }

  reset() {
    this.listReq = {};
    this.fetchList();
  }

  onRowClick(row: Note) {
    this.dialog
      .open(EditNoteDialogComponent, {
        width: "800px",
        data: {
          note: { ...row },
        },
      })
      .afterClosed()
      .subscribe(() => {
        this.fetchList();
      });
  }

  pageChanged(evt: Paging) {
    this.fetchList();
  }
}
