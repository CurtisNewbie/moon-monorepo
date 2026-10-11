import { HttpClient } from "@angular/common/http";
import { Injectable, OnDestroy } from "@angular/core";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Observable, of, Subject } from "rxjs";
import { map, catchError, debounceTime, switchMap, tap } from "rxjs/operators";

export interface RecordBrowseHistoryReq {
  fileKey?: string;
}

interface DirPageUpdate {
  dirKey: string;
  fileKey: string;
}

const DIR_LAST_PAGE_ENDPOINT = '/vfm/history/dir/last-page';

@Injectable({
  providedIn: "root",
})
export class BrowseHistoryRecorder implements OnDestroy {
  // Dir-page writes funnel through a single serialized pipeline: debounceTime
  // coalesces bursts, switchMap keeps at most one request in flight, so a slow
  // earlier write can never land after (and overwrite) a later one.
  private dirPageUpdates = new Subject<DirPageUpdate>();
  // The latest requested pair. Kept until unload (not cleared on dispatch) so a
  // tab close can re-beacon it even while an older write is still in flight.
  private lastDirPage: DirPageUpdate | null = null;
  private beforeUnloadHandler = () => this.flushLastDirPage();

  constructor(private snackBar: MatSnackBar, private http: HttpClient) {
    this.dirPageUpdates
      .pipe(
        debounceTime(300),
        switchMap(({ dirKey, fileKey }) => this.postDirPage(dirKey, fileKey))
      )
      .subscribe();

    window.addEventListener('beforeunload', this.beforeUnloadHandler);
  }

  ngOnDestroy() {
    window.removeEventListener('beforeunload', this.beforeUnloadHandler);
  }

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

  recordDirPage(dirKey: string, fileKey: string) {
    this.lastDirPage = { dirKey, fileKey };
    this.dirPageUpdates.next({ dirKey, fileKey });
  }

  private postDirPage(dirKey: string, fileKey: string): Observable<any> {
    return this.http.post<any>(DIR_LAST_PAGE_ENDPOINT, { dirKey, fileKey }).pipe(
      tap((resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, 'ok', { duration: 6000 });
        }
      }),
      catchError((err) => {
        console.log(err);
        this.snackBar.open('Request failed, unknown error', 'ok', {
          duration: 3000,
        });
        return of(null);
      })
    );
  }

  private flushLastDirPage() {
    const latest = this.lastDirPage;
    if (!latest) {
      return;
    }
    const payload = JSON.stringify({ dirKey: latest.dirKey, fileKey: latest.fileKey });
    navigator.sendBeacon(DIR_LAST_PAGE_ENDPOINT, new Blob([payload], { type: 'application/json' }));
  }

  getDirPage(dirKey: string): Observable<string> {
    return this.http.get<any>(`/vfm/history/dir/last-page?dirKey=${encodeURIComponent(dirKey)}`).pipe(
      map(resp => {
        if (resp.error) return '';
        return resp.data?.fileKey || '';
      }),
      catchError(() => of(''))
    );
  }
}
