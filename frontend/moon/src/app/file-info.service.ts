import { HttpClient, HttpEvent, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { UploadFileParam } from "src/common/file-info";
import { Resp } from "src/common/resp";
import { MatSnackBar } from "@angular/material/snack-bar";

export enum TokenType {
  DOWNLOAD = "DOWNLOAD",
  STREAMING = "STREAMING",
}

@Injectable({
  providedIn: "root",
})
export class FileInfoService {
  constructor(private http: HttpClient, private snackBar: MatSnackBar) {}

  public uploadToMiniFstore(
    uploadParam: UploadFileParam
  ): Observable<HttpEvent<any>> {
    let headers = new HttpHeaders().append(
      "fileName",
      encodeURI(uploadParam.fileName)
    );

    return this.http.put<HttpEvent<any>>("fstore/file", uploadParam.files[0], {
      observe: "events",
      reportProgress: true,
      withCredentials: true,
      headers: headers,
    });
  }

  public generateFileTempToken(
    fileKey: string,
    tokenType: TokenType = TokenType.DOWNLOAD
  ): Observable<Resp<string>> {
    return this.http.post<Resp<string>>(`vfm/open/api/file/token/generate`, {
      fileKey: fileKey,
      tokenType: tokenType,
    });
  }

  public jumpToDownloadUrl(fileKey: string): void {
    this.generateFileTempToken(fileKey).subscribe({
      next: (resp) => {
        if (resp.error) {
          this.snackBar.open(resp.msg, "ok", { duration: 6000 });
          return;
        }
        const token = resp.data;
        const url = "fstore/file/raw?key=" + encodeURIComponent(token);
        window.open(url, "_parent");
      },
    });
  }
}
