import { Injectable } from "@angular/core";
import {
  HttpInterceptor,
  HttpEvent,
  HttpResponse,
  HttpRequest,
  HttpHandler,
  HttpErrorResponse,
  HttpHeaderResponse,
} from "@angular/common/http";
import { Observable } from "rxjs";
import { catchError, filter } from "rxjs/operators";
import { Resp } from "src/common/resp";
import { Router } from "@angular/router";
import { Toaster } from "../notification.service";

/**
 * Intercept http response with 'Resp' as body
 */
@Injectable()
export class RespInterceptor implements HttpInterceptor {
  constructor(private router: Router, private notifi: Toaster) {}

  intercept(
    httpRequest: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    return next.handle(httpRequest).pipe(
      filter((e, i) => {
        if (!(e instanceof HttpResponse || e instanceof HttpHeaderResponse)) {
          return true;
        }

        // console.log("Intercept HttpResponse:", e);

        if (e instanceof HttpResponse) {
          let r: Resp<any> = e.body as Resp<any>;
          if (r.error) {
            this.notifi.toast(r.msg, 6000);
            return false; // filter out this value
          }
        }
        return true;
      })
    );
  }
}
