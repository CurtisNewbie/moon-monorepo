import { Injectable } from "@angular/core";
import {
  HttpInterceptor,
  HttpEvent,
  HttpRequest,
  HttpHandler,
} from "@angular/common/http";
import { Observable } from "rxjs";
import { getToken } from "src/common/api-util";

@Injectable()
export class TokenInterceptor implements HttpInterceptor {

  constructor() { }

  intercept(
    httpRequest: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    let tk = getToken()
    if (tk) {
      httpRequest = httpRequest.clone({
        headers: httpRequest.headers.set("Authorization", tk)
      })
    }
    return next.handle(httpRequest);
  }
}
