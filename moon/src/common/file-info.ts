import { Paging } from "./paging";
import { Option } from "./select-util";

export interface FileInfo {
  id: number;
  uuid: string;
  name: string;
  uploaderName: string;
  uploadTime: Date;
  sizeInBytes: number;
  fileType: FileType;
  updateTime: Date;
  sensitiveMode: string;
  thumbnailToken: string;

  /*
    ---------------------------

    Used by frontend only
    ---------------------------
  */

  /** Label for File Type */
  fileTypeLabel: string;

  /** Label for size */
  sizeLabel: string;

  /**
   * whether file is selected
   */
  _selected: boolean;

  /**
   * whether fileType == 'FILE'
   */
  isFile: boolean;

  /**
   * whether fileType == 'DIR'
   */
  isDir: boolean;

  parentFileName?: string;

  isDisplayable: boolean;

  thumbnailUrl: string;
}

export enum FileType {
  /** File */
  FILE = "FILE",
  /** Directory */
  DIR = "DIR"
}

export function getFileTypeOpts(includesAll: boolean = true): Option<FileType>[] {
  let l = [];
  if (includesAll) l.push({ name: "All", value: null });

  l.push({ name: "File", value: FileType.FILE });
  l.push({ name: "Directory", value: FileType.DIR });
  return l;
}


/** Brief info for DIR type file */
export interface DirBrief {
  id: number;
  uuid: string;
  name: string;
}

/** Parameters used for fetching list of file info */
export interface SearchFileInfoParam {
  /** filename */
  name?: string;
  /** folder no */
  folderNo?: string;
  /** parent file UUID */
  parentFile?: string;
  /** fileType */
  fileType?: FileType;
}

/** Parameters for uploading a file */
export interface UploadFileParam {
  /** name of the file */
  fileName?: string;
  /** file */
  files?: File[];
  /** parent file uuid */
  parentFile?: string;
  /** ignore on duplicate name */
  ignoreOnDupName?: boolean;
}

/** Parameters for fetching list of file info */
export interface FetchFileInfoListParam {
  /** filename */
  filename?: string;
  /** paging  */
  paging?: Paging;
  /** folder no */
  folderNo?: string;
  /** parent file UUID */
  parentFile?: string;
}

/**
 * Empty object with all properties being null values
 */
export function emptyUploadFileParam(): UploadFileParam {
  return {
    files: [],
    fileName: null,
  };
}

export interface FileAccessGranted {
  /** id of this file_sharing record */
  id: number;
  /** id of user */
  userId?: number;
  /* userNo */
  userNo?: string;
  /** user who is granted access to this file*/
  username: string;
  /** the date that this access is granted */
  createDate: Date;
  /** the access is granted by */
  createdBy: string;
}