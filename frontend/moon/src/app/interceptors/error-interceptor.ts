import { Injectable } from "@angular/core";
import {
  HttpInterceptor,
  HttpEvent,
  HttpRequest,
  HttpHandler,
  HttpErrorResponse,
} from "@angular/common/http";
import { Observable, throwError } from "rxjs";
import { catchError } from "rxjs/operators";
import { UserService } from "../user.service";
import { Resp } from "src/common/resp";
import { Router } from "@angular/router";
import { MatSnackBar } from "@angular/material/snack-bar";

/**
 * Intercept http error response
 */
@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  constructor(
    private userService: UserService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {}

  intercept(
    httpRequest: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(httpRequest).pipe(
      catchError((e) => {
        if (e instanceof HttpErrorResponse) {
          // console.log("Http error response: ", e);

          if (e.status === 401) {
            this.snackBar.open("Please login first", "ok", { duration: 3000 });
            if (this.router.url != "/register") {
              this.userService.logout();
            }
          } else if (e.status === 403) {
            let r: Resp<any> = e.error as Resp<any>;
            if (r) {
              this.snackBar.open(r.msg, "ok", { duration: 6000 });
            } else {
              this.snackBar.open("Forbidden", "ok", { duration: 1000 });
            }
          } else {
            this.snackBar.open(
              "Unknown server error, please try again later",
              "ok",
              { duration: 3000 }
            );
          }
          return throwError(e);
        }
      })
    );
  }
}
