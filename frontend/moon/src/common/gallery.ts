import { Paging } from "./paging";

export interface Gallery {
  id: string;
  galleryNo: string;
  userNo: string;
  name: string;
  createTime: string;
  createBy: string;
  updateTime: string;
  updateBy: string;
  isOwner: boolean;
}

export interface ListGalleryImagesResp {
  images: { thumbnailToken: string, fileTempToken: string }[];
  paging: Paging;
}

export interface GalleryBrief {
  /** gallery no */
  galleryNo: string;
  /** name of the gallery */
  name: string;
}