import { Pipe, PipeTransform } from "@angular/core";
import { I18n } from "./i18n.service";

@Pipe({
  name: "trl",
  pure: true,
})
export class TrlPipe implements PipeTransform {
  constructor(private i18n: I18n) {}

  transform(mod: string, key: string, ...args: any[]): string {
    return this.i18n.trl(mod, key, ...args);
  }
}