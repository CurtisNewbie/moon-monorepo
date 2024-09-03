import { Paging } from "./paging";

export interface UserInfo {
  id: number;
  username: string;
  role: string; // deprecated

  /** whether the user is disabled, 0-normal, 1-disabled */
  isDisabled: number;

  /** review status */
  reviewStatus: string;

  roleNo: string;
  userNo: string;
  roleName: string;
  createTime: Date;
  updateTime: Date;
  updateBy: string;
  createBy: string;
}

export interface FetchUserInfoResp {
  list: UserInfo[];
  paging: Paging;
}

export enum UserIsDisabledEnum {
  /**
   * User is in normal state
   */
  NORMAL = 0,

  /**
   * User is disabled
   */
  IS_DISABLED = 1,
}

export interface UserIsDisabledOption {
  name: string;
  value: number;
}

export const USER_IS_DISABLED_OPTIONS: UserIsDisabledOption[] = [
  { name: "normal", value: UserIsDisabledEnum.NORMAL },
  { name: "disabled", value: UserIsDisabledEnum.IS_DISABLED },
];

/**
 * Parameters for adding a new user
 */
export interface AddUserParam {
  /** username */
  username: string;
  /** password */
  password: string;
  /** user role */
  userRole: string;
}

/**
 * Parameters for registration
 */
export interface RegisterUserParam {
  /** username */
  username: string;
  /** password */
  password: string;
}

/**
 * Parameters for changing password
 */
export interface ChangePasswordParam {
  /**
   * Previous password
   */
  prevPassword: string;

  /**
   * New password
   */
  newPassword: string;
}

/**
 * Empty object with all properties being null values
 */
export function emptyChangePasswordParam(): ChangePasswordParam {
  return {
    prevPassword: null,
    newPassword: null,
  };
}

/**
 * Parameters for fetching user info
 */
export interface FetchUserInfoParam {
  username?: string;
  isDisabled?: UserIsDisabledEnum | number;
  roleNo?: string;
  paging?: Paging;
}

export interface UpdateUserInfoParam {
  /**
   * User id
   */
  id: number;

  /**
   * User role
   */
  role: string;

  /**
   * User's is_disabled status
   */
  isDisabled: number;
}
