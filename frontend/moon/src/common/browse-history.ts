import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";

export interface RecordBrowseHistoryReq {
  fileKey?: string;
}

@Injectable({
  providedIn: "root",
})
export class BrowseHistoryRecorder {
  constructor(private snackBar: MatSnackBar, private http: HttpClient) {}

  record(fileKey: string) {
    let req: RecordBrowseHistoryReq = { fileKey: fileKey };
    this.http.post<any>(`/vfm/history/record-browse-history`, req).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
      },
      error: (err) => {
        console.log(err);
        this.snackBar.open("Request failed, unknown error", "ok", {
          duration: 3000,
        });
      },
    });
  }
}
