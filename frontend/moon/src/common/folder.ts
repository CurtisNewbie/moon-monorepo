
/** Virtual Folder */
export interface VFolder {
  id: number;

  /** when the record is created */
  createTime: Date;

  /** who created this record */
  createBy: string;

  /** when the record is updated */
  updateTime: Date;

  /** who updated this record */
  updateBy: string;

  /** folder no */
  folderNo: string;

  /** name of the folder */
  name: string;

  /** ownership */
  ownership: string;
}

export interface VFolderBrief {
  /** folder no */
  folderNo: string;
  /** name of the folder */
  name: string;
}
