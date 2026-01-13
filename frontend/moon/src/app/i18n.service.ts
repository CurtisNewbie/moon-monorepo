import { Injectable } from "@angular/core";
import { enUSDict, zhCNDict } from "./translate.dict";

@Injectable({
  providedIn: "root",
})
export class I18n {
  public options: string[] = ["zh-CN", "en-US"];
  public optionsLabels = [
    { value: "zh-CN", name: "中文" },
    { value: "en-US", name: "Eng" },
  ];
  public curr: string = "en-US";
  private dicts = {
    "zh-CN": zhCNDict,
    "en-US": enUSDict,
  };
  constructor() {
    let prev = localStorage.getItem("i18n.lang");
    if (prev) {
      this.change(prev);
    }
  }

  change(op) {
    for (let v of this.options) {
      if (op === v) {
        this.curr = op;
        localStorage.setItem("i18n.lang", op);
        return;
      }
    }
  }

  trl(mod: string, name: string, ...args: any) {
    let p: string;
    if (mod == "") {
      p = this.dicts[this.curr][name] ?? name;
    } else {
      p = this.dicts[this.curr][mod][name] ?? name;
    }
    if (args && args.length > 0) {
      return this.NamedSprintfkv(p, ...args);
    }
    return p;
  }

  // Regular expression to match named placeholders like ${name}
  namedFmtPat = /\${[a-zA-Z0-9\/\-_\. ]+}/g;

  // Format message using key-value pairs, e.g., '${startTime} ${message}'
  //
  // e.g., NamedSprintfkv("my name is ${name}", "name", "slim shady")
  NamedSprintfkv<T extends string | number | boolean | bigint>(
    pat: string,
    ...kv: T[]
  ): string {
    if (kv.length < 1) {
      return pat;
    }

    const p: Record<string, any> = {};
    let lastKey: string | null = null;

    for (let i = 0; i < kv.length; i++) {
      const item = kv[i];
      const strValue = String(item);

      if (i % 2 === 0) {
        lastKey = strValue;
      } else if (lastKey !== null) {
        p[lastKey] = strValue;
        lastKey = null;
      }
    }

    return pat.replace(this.namedFmtPat, (match: string) => {
      const key = match.slice(2, -1); // Remove '${' and '}'
      const value = p[key];
      if (value === undefined || value === null) {
        return match; // Keep the original placeholder if key not found
      }
      return String(value);
    });
  }
}
