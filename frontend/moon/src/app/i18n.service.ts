import { Injectable } from "@angular/core";
import { enUSDict, zhCNDict } from "./translate.dict";

@Injectable({
  providedIn: "root",
})
export class I18n {
  public options: string[] = ["zh-CN", "en-US"];
  public optionsLabels = [
    { value: "zh-CN", name: "中文" },
    { value: "en-US", name: "English" },
  ];
  public curr: string = "en-US";
  private dicts = {
    "zh-CN": zhCNDict,
    "en-US": enUSDict,
  };
  constructor() {}

  change(op) {
    for (let v of this.options) {
      if (op === v) {
        this.curr = op;
        return;
      }
    }
  }

  trl(mod: string, name: string) {
    if (mod == "") {
      return this.dicts[this.curr][name] ?? name;
    }
    return this.dicts[this.curr][mod][name] ?? name;
  }
}
