import { Component, Inject, OnInit } from "@angular/core";
import { Note } from "../mng-note/mng-note.component";
import { HttpClient } from "@angular/common/http";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";
import { ConfirmDialog } from "src/common/dialog";

export interface DialogData {
  note: Note;
}

export interface UpdateNoteReq {
  recordId?: string; // Required.
  title?: string;
  content?: string;
}

export interface ApiDeleteNoteReq {
  recordId?: string;
}

@Component({
  selector: "app-edit-note-dialog",
  templateUrl: "./edit-note-dialog.component.html",
  styleUrls: ["./edit-note-dialog.component.css"],
})
export class EditNoteDialogComponent implements OnInit {
  edit = false;

  constructor(
    private http: HttpClient,
    private snackBar: MatSnackBar,
    private confirmDialog: ConfirmDialog,
    public dialogRef: MatDialogRef<EditNoteDialogComponent, DialogData>,
    @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {}

  ngOnInit(): void {
    this.edit = false;
  }

  switchEdit(): void {
    this.edit = !this.edit;
  }

  updateNote() {
    let req: UpdateNoteReq = {
      recordId: this.data.note.recordId,
      title: this.data.note.title,
      content: this.data.note.content,
    };
    this.http
      .post<any>(`/user-vault/open/api/note/update-note`, req)
      .subscribe({
        next: (resp) => {
          if (resp.error) {
            this.snackBar.open(resp.msg, "ok", { duration: 6000 });
            return;
          }
          this.dialogRef.close();
        },
        error: (err) => {
          console.log(err);
          this.snackBar.open("Request failed, unknown error", "ok", {
            duration: 3000,
          });
        },
      });
  }

  deleteNote() {
    this.confirmDialog.show(
      `Delete Note`,
      [`Delete '${this.data.note.title}'?`],
      () => {
        let req: ApiDeleteNoteReq = {
          recordId: this.data.note.recordId,
        };
        this.http
          .post<any>(`/user-vault/open/api/note/delete-note`, req)
          .subscribe({
            next: (resp) => {
              if (resp.error) {
                this.snackBar.open(resp.msg, "ok", { duration: 6000 });
                return;
              }
              this.dialogRef.close();
            },
            error: (err) => {
              console.log(err);
              this.snackBar.open("Request failed, unknown error", "ok", {
                duration: 3000,
              });
            },
          });
      }
    );
  }
}
