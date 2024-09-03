import { HttpClient } from "@angular/common/http";
import { Component, OnInit, ViewChild } from "@angular/core";
import { ActivatedRoute, ParamMap } from "@angular/router";
import { PdfJsViewerComponent } from "ng2-pdfjs-viewer";

@Component({
  selector: "app-viewer",
  templateUrl: "./pdf-viewer.component.html",
  styleUrls: ["./pdf-viewer.component.css"],
})
export class PdfViewerComponent implements OnInit {
  uuid: string;
  name: string;

  @ViewChild("pdfViewer", { static: true })
  pdfViewer: PdfJsViewerComponent;

  constructor(private route: ActivatedRoute, private httpClient: HttpClient) { }

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      this.name = params.get("name");
      this.uuid = params.get("uuid");
      this.httpClient
        .get(params.get("url"), {
          responseType: "blob",
          observe: "body",
        })
        .subscribe({
          next: (blob) => {
            this.pdfViewer.pdfSrc = blob;
            this.pdfViewer.refresh();

            // jump to the previous page when document is loaded
            this.pdfViewer.onDocumentLoad.subscribe(() => {
              let page = localStorage.getItem(this.storageKey(this.uuid));
              if (page) {
                this.pdfViewer.page = parseInt(page);
              }
            });

            // record the last page
            this.pdfViewer.onPageChange.subscribe((page) => {
              localStorage.setItem(this.storageKey(this.uuid), page);
            });
          },
        });
    });
  }

  private storageKey(uuid): string {
    return `pdf:viewer:${uuid}`;
  }
}
