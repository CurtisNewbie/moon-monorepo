import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-txt-viewer',
  templateUrl: './txt-viewer.component.html',
  styleUrls: ['./txt-viewer.component.css']
})
export class TxtViewerComponent implements OnInit {

  uuid: string;
  name: string;
  content: string;

  constructor(private route: ActivatedRoute, private httpClient: HttpClient) { }

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      this.name = params.get("name");
      this.uuid = params.get("uuid");
      this.httpClient
        .get(params.get("url"), { observe: "body", responseType: "text" })
        .subscribe({
          next: (txt) => {
            this.content = txt; 
          }
        });
    });
  }

}
