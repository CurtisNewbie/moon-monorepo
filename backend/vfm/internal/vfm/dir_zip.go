package vfm

import (
	"fmt"
	"strings"
	"unicode"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

type GenerateDirZipReq struct {
	FileKey string `json:"fileKey" valid:"notEmpty"`
}

// dirZipChild is the projection of one direct file child that goes into the archive.
type dirZipChild struct {
	Name         string
	FstoreFileId string
}

// GenDirBundleKey mints a bundle key that downloads the directory as a flat zip
// archive of its direct file children.
//
// Only the caller's own, live directory is accepted. Sub-directories and their
// contents are excluded by construction: the query selects direct file rows of
// this one directory. A child without an fstore file id cannot be served, so it
// is skipped with a warning instead of failing the whole archive.
func GenDirBundleKey(rail miso.Rail, db *gorm.DB, req GenerateDirZipReq, user flow.User) (string, error) {
	dir, ok, err := findFile(rail, db, req.FileKey)
	if err != nil {
		return "", errs.Wrapf(err, "failed to find directory, fileKey: %v", req.FileKey)
	}
	if !ok {
		return "", ErrFileNotFound.New()
	}
	// ownership comes before the deleted/type checks so a stranger learns
	// nothing about another user's uuid beyond its existence
	if dir.UploaderNo != user.UserNo {
		return "", errs.ErrNotPermitted.New()
	}
	if dir.IsLogicDeleted == DelY {
		return "", ErrFileDeleted.New()
	}
	if dir.FileType != FileTypeDir {
		return "", errs.NewErrf("Only a directory can be downloaded as a zip archive")
	}

	var children []dirZipChild
	if err := dbquery.NewQuery(rail, db).
		Select("name, fstore_file_id").
		Table("file_info").
		Eq("parent_file", req.FileKey).
		Eq("file_type", FileTypeFile).
		Eq("uploader_no", user.UserNo).
		Eq("is_logic_deleted", DelN).
		Eq("is_del", 0).
		Eq("hidden", 0).
		Where("fstore_file_id is not null and fstore_file_id != ''").
		Order("name asc, id asc").
		ScanVal(&children); err != nil {
		return "", errs.Wrapf(err, "failed to query files in directory, fileKey: %v", req.FileKey)
	}

	items := make([]fstore.BundleItem, 0, len(children))
	used := make(map[string]struct{}, len(children))
	for _, c := range children {
		if c.FstoreFileId == "" {
			rail.Warnf("Skipping zip entry, file has no mini-fstore file_id, dir: %v, name: %v", req.FileKey, c.Name)
			continue
		}
		if problem := dirZipNameProblem(c.Name); problem != "" {
			return "", errs.NewErrf("File %q in directory %q cannot be included in the zip archive: %s", c.Name, dir.Name, problem)
		}
		items = append(items, fstore.BundleItem{
			FileId: c.FstoreFileId,
			Name:   dirZipEntryName(c.Name, used),
		})
	}

	key, err := fstore.GenBundleKey(rail, dir.Name, items)
	if err != nil {
		return "", errs.Wrapf(err, "failed to generate bundle key for directory, fileKey: %v", req.FileKey)
	}
	rail.Infof("Generated bundle key for directory %v with %d entries", req.FileKey, len(items))
	return key, nil
}

// dirZipNameProblem reports why a file name cannot be used verbatim as an entry
// name inside the flat archive, or "" when it can.
//
// A name that would nest or escape (a path separator, a ".." segment), that
// cannot be represented reliably (a control character), or that a zip tool would
// silently alter (an empty name, a trailing dot or space) is rejected instead of
// being rewritten, so the entry name always equals the file name the user sees.
func dirZipNameProblem(name string) string {
	if name == "" {
		return "the name is empty"
	}
	if strings.ContainsAny(name, `/\`) {
		return "the name contains a path separator"
	}
	for _, seg := range strings.FieldsFunc(name, func(r rune) bool { return r == '/' || r == '\\' }) {
		if seg == ".." {
			return "the name contains a '..' path segment"
		}
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return "the name contains a control character"
		}
	}
	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return "the name ends with a dot or a space"
	}
	return ""
}

// dirZipEntryName resolves one entry name, suffixing a name that is already
// taken with " (n)" before its extension. Names are assigned in the fixed query
// order, so the same directory always produces the same entry names.
func dirZipEntryName(name string, used map[string]struct{}) string {
	if _, exists := used[name]; !exists {
		used[name] = struct{}{}
		return name
	}

	base, ext := splitFileName(name)
	for n := 1; ; n++ {
		candidate := fmt.Sprintf("%s (%d)%s", base, n, ext)
		if _, exists := used[candidate]; !exists {
			used[candidate] = struct{}{}
			return candidate
		}
	}
}

// splitFileName splits a filename into the part before the final dot and the
// extension suffix. A name without an extension, or a dotfile whose only dot is
// the leading one, has no extension.
func splitFileName(name string) (string, string) {
	i := strings.LastIndex(name, ".")
	if i <= 0 {
		return name, ""
	}
	return name[:i], name[i:]
}
