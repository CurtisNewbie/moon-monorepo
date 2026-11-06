/** One node of a Zip archive's tree, as the preview endpoint returns it. */
export interface ZipTreeNode {
  /** one path segment; empty for the response's root node only */
  name: string;
  /** true for a folder, false for a file */
  dir: boolean;
  /** uncompressed bytes; files only */
  size?: number;
  /** ordered children; folders only */
  children?: ZipTreeNode[];
}

/** A Zip archive's tree, plus how far the server read into the archive. */
export interface ZipTreePreview {
  /** the archive's root folder node, with an empty name */
  tree: ZipTreeNode;
  /** true when a cap left entries out */
  truncated: boolean;
  /** zero or more of depth and entries */
  truncatedReasons: string[];
  /** how many tree nodes the tree represents */
  includedEntries: number;
}

/** Truncation cause: entries deeper than the depth cap were left out. */
export const ZIP_TREE_TRUNCATED_BY_DEPTH = 'depth';

/** Truncation cause: entries past the node cap were left out. */
export const ZIP_TREE_TRUNCATED_BY_ENTRIES = 'entries';
