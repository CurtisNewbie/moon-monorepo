import { FileInfo, FileType } from "./file-info";

const KB_UNIT: number = 1024;
const MB_UNIT: number = 1024 * 1024;
const GB_UNIT: number = 1024 * 1024 * 1024;

const videoSuffix = new Set(["mp4", "mov", "webm", "ogg"]);
const imageSuffix = new Set([
  "jpeg",
  "jpg",
  "gif",
  "png",
  "svg",
  "bmp",
  "webp",
  "apng",
  "avif",
]);
const textSuffix = new Set([
  "conf",
  "txt",
  "yml",
  "yaml",
  "properties",
  "json",
  "sh",
  "md",
  "java",
  "js",
  "ts",
  "css",
  "list",
  "service",
  "xml",
  "sql",
  "go",
]);
const webpageSuffix = new Set(["html"]);

const suffixIcon: [Set<string>, string][] = [
  [new Set(["pdf"]), "bi-filetype-pdf"],
  [new Set(["zip", "7z", "jar"]), "bi-file-zip"],
  [new Set(["yml", "yaml"]), "bi-filetype-yml"],
  [
    new Set([
      "txt",
      "conf",
      "properties",
      "json",
      "list",
      "doc",
      "docx",
      "service",
      "md",
      "conf",
    ]),
    "bi-file-text",
  ],
  [new Set(["java"]), "bi-filetype-java"],
  [new Set(["js", "ts"]), "bi-filetype-js"],
  [new Set(["html"]), "bi-filetype-html"],
  [new Set(["css"]), "bi-filetype-css"],
  [new Set(["sh"]), "bi-filetype-sh"],
  [new Set(["csv"]), "bi-filetype-csv"],
  [new Set(["xls"]), "bi-filetype-xls"],
  [new Set(["xlsx"]), "bi-filetype-xlss"],
  [new Set(["doc"]), "bi-filetype-doc"],
  [new Set(["xml"]), "bi-filetype-xml"],
  [new Set(["py"]), "bi-filetype-py"],
  [new Set(["json"]), "bi-filetype-json"],
  [new Set(["exe"]), "bi-filetype-exe"],
  [new Set(["dmg", "iso"]), "bi-hdd"],
];

export function resolveSize(sizeInBytes: number): string {
  if (sizeInBytes > GB_UNIT) {
    return divideUnit(sizeInBytes, GB_UNIT) + " gb";
  }
  if (sizeInBytes > MB_UNIT) {
    return divideUnit(sizeInBytes, MB_UNIT) + " mb";
  }
  return divideUnit(sizeInBytes, KB_UNIT) + " kb";
}

export function divideUnit(size: number, unit: number): string {
  return (size / unit).toFixed(1);
}

export function suffix(name: string): string {
  let i = name.lastIndexOf(".");
  if (i < 0 || i == name.length - 1) return "";

  let suffix = name.slice(i + 1);
  return suffix.toLowerCase();
}

export function guessFileIconClz(f: any): string {
  if (f.fileType == FileType.DIR) {
    return "bi-folder-fill";
  }
  let suf = suffix(f.name);
  if (!suf) {
    return "bi-file-binary";
  }
  for (let u of suffixIcon) {
    if (u[0].has(suf)) {
      return u[1];
    }
  }
  return "bi-file-text";
}

export function isWebpage(fname: string): boolean {
  return fileSuffixAnyMatch(fname, webpageSuffix);
}

export function isTxt(fname: string): boolean {
  return fileSuffixAnyMatch(fname, textSuffix);
}
export function fileSuffixAnyMatch(
  name: string,
  candidates: Set<string>
): boolean {
  let i = name.lastIndexOf(".");
  if (i < 0 || i == name.length - 1) return false;

  let suffix = name.slice(i + 1);
  return candidates.has(suffix.toLowerCase());
}

export function isPdf(filename: string): boolean {
  return filename.toLowerCase().indexOf(".pdf") != -1;
}

export function isStreamableVideo(filename: string): boolean {
  return fileSuffixAnyMatch(filename, videoSuffix);
}

export function isImageByName(filename: string): boolean {
  return fileSuffixAnyMatch(filename, imageSuffix);
}

export function isImage(f: FileInfo): boolean {
  if (f == null || !f.isFile) return false;
  return isImageByName(f.name);
}

export function canPreview(filename: string): boolean {
  return (
    isPdf(filename) ||
    isImageByName(filename) ||
    isStreamableVideo(filename) ||
    isTxt(filename)
  );
}
