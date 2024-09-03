import { Component, OnInit } from "@angular/core";
import { PagingController } from "src/common/paging";
import {
  animateElementExpanding,
  getExpanded,
  isIdEqual,
} from "src/animate/animate-util";
import { Toaster } from "../notification.service";
import { environment } from "src/environments/environment";
import { isEnterKey } from "src/common/condition";
import { copyToClipboard } from "src/common/clipboard";
import { HttpClient } from "@angular/common/http";

export interface UserToken {

    id: number;

    /** secret key */
    secretKey: string;

    /** name of the key */
    name: string;

    /** when the key is expired */
    expirationTime: Date;

    /** when the record is created */
    createTime: Date;

}

@Component({
  selector: "app-manage-keys",
  templateUrl: "./manage-keys.component.html",
  styleUrls: ["./manage-keys.component.css"],
  animations: [animateElementExpanding()],
})
export class ManageKeysComponent implements OnInit {
  readonly columns: string[] = [
    "id",
    "name",
    "secretKey",
    "expirationTime",
    "createTime",
  ];
  expandedElement: UserToken = null;
  tokens: UserToken[] = [];
  pagingController: PagingController;
  query = {
    name: "",
  };
  panelDisplayed: boolean = false;
  password: string = null;
  newUserKeyName: string = null;

  idEquals = isIdEqual;
  getExpandedEle = (row) => getExpanded(row, this.expandedElement);
  isEnter = isEnterKey;
  copyToClipboard = copyToClipboard;

  constructor(private http: HttpClient, private toaster: Toaster) {}

  ngOnInit() {}

  mask(k: string): string {
    return k.length > 0
      ? k.substring(0, 5) + "*********" + k.substring(k.length - 5)
      : "";
  }

  fetchList() {
    this.http
      .post<any>(`${environment.uservault}/open/api/user/key/list`, {
        payload: { name: this.query.name },
        paging: this.pagingController.paging,
      })
      .subscribe((resp) => {
        if (resp.data) {
          this.tokens = [];
          if (resp.data.payload) {
            for (let r of resp.data.payload) {
              if (r.expirationTime)
                r.expirationTime = new Date(r.expirationTime);
              if (r.createTime) r.createTime = new Date(r.createTime);
              this.tokens.push(r);
            }
          }
          this.pagingController.onTotalChanged(resp.data.paging);
          if (this.panelDisplayed) this.panelDisplayed = false;
        }
      });
  }

  reset() {
    this.expandedElement = null;
    this.pagingController.firstPage();
    this.query = {
      name: "",
    };
  }

  generateRandomKey() {
    if (!this.password) {
      this.toaster.toast("Please enter password");
      return;
    }
    if (!this.newUserKeyName) {
      this.toaster.toast("Please enter key name");
      return;
    }

    const pw = this.password;
    const keyName = this.newUserKeyName;

    this.password = null;

    this.http
      .post<void>(`${environment.uservault}/open/api/user/key/generate`, {
        password: pw,
        keyName: keyName,
      })
      .subscribe({
        next: (resp) => {
          this.fetchList();
          this.newUserKeyName = null;
          this.panelDisplayed = false;
        },
      });
  }

  deleteUserKey(id: number) {
    this.http
      .post<void>(`${environment.uservault}/open/api/user/key/delete`, {
        userKeyId: id,
      })
      .subscribe({
        complete: () => this.fetchList(),
      });
  }

  togglePanel() {
    this.panelDisplayed = !this.panelDisplayed;
    this.password = null;
  }

  onPagingControllerReady(pc) {
    this.pagingController = pc;
    this.pagingController.onPageChanged = () => this.fetchList();
    this.fetchList();
  }
}
