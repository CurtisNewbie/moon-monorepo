package api

const (
	// DefaultZipTreeMaxDepth is the fallback depth limit applied when a preview
	// request carries no limit (zero or negative).
	DefaultZipTreeMaxDepth = 32

	// DefaultZipTreeMaxEntries is the fallback node limit applied when a
	// preview request carries no limit (zero or negative).
	DefaultZipTreeMaxEntries = 10000
)

const (
	// ZipTreeTruncatedByDepth marks a preview whose depth limit left entries out.
	ZipTreeTruncatedByDepth = "depth"

	// ZipTreeTruncatedByEntries marks a preview whose node limit left entries out.
	ZipTreeTruncatedByEntries = "entries"
)

type ZipTreeReq struct {
	// zip file's mini-fstore file_id.
	FileId string `json:"fileId" desc:"file_id of the zip file"`

	// Maximum number of path segments an entry may sit at, counted from the
	// archive root. Zero means the server default.
	MaxDepth int `json:"maxDepth" desc:"maximum path segment depth included, zero means the server default"`

	// Maximum number of tree nodes represented, counting files, explicit
	// folders and inferred folders. Zero means the server default.
	MaxEntries int `json:"maxEntries" desc:"maximum tree nodes represented, zero means the server default"`
}

// ZipTreeNode is one node of a zip file's tree. A folder node has children; a
// file node has an uncompressed size. An empty folder omits children, and a
// zero-byte file omits size.
type ZipTreeNode struct {
	Name     string         `json:"name" desc:"one path segment, empty for the root node only"`
	Dir      bool           `json:"dir" desc:"true for a folder node, false for a file node"`
	Size     int64          `json:"size,omitempty" desc:"uncompressed size in bytes, files only"`
	Children []*ZipTreeNode `json:"children,omitempty" desc:"ordered children, folders only"`
}

type ZipTreeResp struct {
	Tree             ZipTreeNode `json:"tree" desc:"root folder node, its name is empty"`
	Truncated        bool        `json:"truncated" desc:"true when a cap left nodes out"`
	TruncatedReasons []string    `json:"truncatedReasons" desc:"zero or more of 'depth' and 'entries', empty when not truncated"`
	IncludedEntries  int         `json:"includedEntries" desc:"number of tree nodes represented in the tree"`
}
