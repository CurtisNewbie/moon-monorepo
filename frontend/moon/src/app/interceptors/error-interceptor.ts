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
import { Toaster } from "../notification.service";
import { Resp } from "src/common/resp";

/**
 * Intercept http error response
 */
@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  constructor(
    private userService: UserService,
    private toaster: Toaster
  ) { }

  intercept(
    httpRequest: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(httpRequest).pipe(
      catchError((e) => {
        if (e instanceof HttpErrorResponse) {
          // console.log("Http error response: ", e);

          if (e.status === 401) {
            this.toaster.toast("Please login first");
            this.userService.logout();
          } else if (e.status === 403) {
            let r: Resp<any> = e.error as Resp<any>;
            if (r) {
              this.toaster.toast(r.msg, 6000);
            } else {
              this.toaster.toast("Forbidden", 1000);
            }
          } else {
            this.toaster.toast("Unknown server error, please try again later");
          }
          return throwError(e);
        }
      })
    );
  }

}
