import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Observable, of } from "rxjs";
import { map, catchError } from "rxjs/operators";

export interface RecordBrowseHistoryReq {
  fileKey?: string;
}

export interface DirLastPageRes {
  page: number;
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

  recordDirPage(dirKey: string, page: number) {
    this.http.post<any>(`/vfm/history/dir/last-page`, { dirKey, page }).subscribe({
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

  getDirPage(dirKey: string): Observable<number> {
    return this.http.get<any>(`/vfm/history/dir/last-page?dirKey=${encodeURIComponent(dirKey)}`).pipe(
      map(resp => {
        if (resp.error) return 1;
        return resp.data?.page || 1;
      }),
      catchError(() => of(1))
    );
  }
}
