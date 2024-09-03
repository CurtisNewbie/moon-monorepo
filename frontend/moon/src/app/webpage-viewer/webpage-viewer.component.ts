import { HttpClient } from "@angular/common/http";
import { Component, OnInit } from "@angular/core";
import { ActivatedRoute, ParamMap } from "@angular/router";
import { DomSanitizer } from "@angular/platform-browser";
import DOMPurify from 'dompurify';

@Component({
  selector: "app-webpage-viewer",
  templateUrl: "./webpage-viewer.component.html",
  styleUrls: ["./webpage-viewer.component.css"],
})
export class WebpageViewerComponent implements OnInit {
  uuid: string;
  name: string;
  url: string;
  content: any;

  constructor(
    private route: ActivatedRoute,
    private httpClient: HttpClient,
    private domSanitizer: DomSanitizer
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      this.name = params.get("name");
      this.uuid = params.get("uuid");
      this.url = params.get("url");
      this.httpClient
        .get(params.get("url"), { observe: "body", responseType: "text" })
        .subscribe({
          next: (txt) => {
            txt = DOMPurify.sanitize(txt, {FORCE_BODY: true});
            this.content = this.domSanitizer.bypassSecurityTrustHtml(txt);
          },
        });
    });
  }
}
