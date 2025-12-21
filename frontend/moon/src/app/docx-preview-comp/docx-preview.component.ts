import { HttpClient } from "@angular/common/http";
import { Component, ElementRef, OnInit, ViewChild } from "@angular/core";
import { ActivatedRoute } from "@angular/router";

export interface DocxPreviewDat {
  url: string;
}

@Component({
  selector: "app-docx-preview",
  templateUrl: "./docx-preview.component.html",
  styleUrls: ["./docx-preview.component.css"],
})
export class DocxPreviewComponent implements OnInit {
  url: string;
  name: string;

  @ViewChild("docxdoc", { static: true })
  docxdoc: ElementRef;

  constructor(private client: HttpClient, private route: ActivatedRoute) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe((r) => {
      this.url = r.get("url");
      this.name = r.get("name");
      if (this.url) {
        this.render();
      }
    });
  }

  render() {
    this.client
      .get(this.url, {
        responseType: "arraybuffer", // Key option to receive binary data
      })
      .subscribe((d) => {
        let docx = (window as any).docx;
        docx
          .renderAsync(d, this.docxdoc.nativeElement)
          .then((x) => console.log("docx: finished"));
      });
  }
}
