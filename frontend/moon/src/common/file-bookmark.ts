import { Injectable } from "@angular/core";

export interface TempFile {
  thumbnailUrl?: string;
  fileType?: string;
  fileKey?: string;
  name?: string;
}

@Injectable({
  providedIn: "root",
})
export class FileBookmark {
  bucket: Map<string, TempFile> = new Map();

  constructor() {}

  add(f: TempFile) {
    this.bucket.set(f.fileKey, f);
  }

  del(fileKey: string) {
    this.bucket.delete(fileKey);
  }

  has(fileKey: string): boolean {
    return this.bucket.has(fileKey);
  }

  clear() {
    this.bucket.clear();
  }
}
