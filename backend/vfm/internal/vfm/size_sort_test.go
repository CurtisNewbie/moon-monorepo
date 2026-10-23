package vfm

import (
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// --- fixtures ---------------------------------------------------------------

// insertSortTestChildSized inserts a file row of the given size under parent and
// returns its uuid.
func insertSortTestChildSized(t *testing.T, db *gorm.DB, parent, name, userNo string, size int64) string {
	t.Helper()
	uuid := "sst-child-" + randutil.ERand(8)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, seq_key)
		VALUES (?,?,?,?,?,?,?,?)`,
		name, uuid, size, FileTypeFile, parent, "", userNo, "").Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// insertSortTestSubDir inserts a directory row of the given size under parent.
func insertSortTestSubDir(t *testing.T, db *gorm.DB, parent, name, userNo string, size int64) string {
	t.Helper()
	uuid := "sst-subdir-" + randutil.ERand(8)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no)
		VALUES (?,?,?,?,?,?,?)`,
		name, uuid, size, FileTypeDir, parent, "", userNo).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// listSortedFiles lists dir with the given request fields (the parent and page
// are filled in) and returns the payload.
func listSortedFiles(t *testing.T, db *gorm.DB, dir string, req ListFileReq, user flow.User) []ListedFile {
	t.Helper()
	parent := dir
	req.ParentFile = &parent
	req.Page = miso.Paging{Page: 1, Limit: 100}
	res, err := ListFiles(miso.EmptyRail(), db, req, user)
	if err != nil {
		t.Fatal(err)
	}
	return res.Payload
}

// listSortedFileNames lists dir with the given request fields and returns the names.
func listSortedFileNames(t *testing.T, db *gorm.DB, dir string, req ListFileReq, user flow.User) []string {
	t.Helper()
	payload := listSortedFiles(t, db, dir, req, user)
	names := make([]string, 0, len(payload))
	for _, f := range payload {
		names = append(names, f.Name)
	}
	return names
}

// --- ordering ---------------------------------------------------------------

func TestListFilesSortBySizeDefaultsToDescending(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// inserted in neither size nor id order
	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size"}, user)
	want := []string{"large.bin", "medium.bin", "small.bin"} // omitted direction → desc
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size desc: got %v, want %v", got, want)
	}
}

func TestListFilesSortBySizeHonorsAscending(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: "asc"}, user)
	want := []string{"small.bin", "medium.bin", "large.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size asc: got %v, want %v", got, want)
	}
}

func TestListFilesSortBySizeBreaksTiesById(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: c=1, a=2, b=3; all the same size
	insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 500)

	cases := []struct {
		orderDir string
		want     []string
	}{
		{"", []string{"b.bin", "a.bin", "c.bin"}},     // omitted → desc, id desc
		{"desc", []string{"b.bin", "a.bin", "c.bin"}}, // id desc
		{"asc", []string{"c.bin", "a.bin", "b.bin"}},  // id asc
	}
	for _, c := range cases {
		got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: c.orderDir}, user)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("equal sizes orderDir=%q: got %v, want %v", c.orderDir, got, c.want)
		}
	}
}

func TestListFilesSortBySizePagesWithoutRepeatOrSkip(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: c=1, a=2, b=3; equal sizes so only the id tie-break orders them
	insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 500)

	parent := dir
	limit := 2
	cases := []struct {
		orderDir string
		want     []string
	}{
		{"desc", []string{"b.bin", "a.bin", "c.bin"}},
		{"asc", []string{"c.bin", "a.bin", "b.bin"}},
	}
	for _, c := range cases {
		got := make([]string, 0, 3)
		for page := 1; page <= 2; page++ {
			res, err := ListFiles(miso.EmptyRail(), db, ListFileReq{
				Page:       miso.Paging{Page: page, Limit: limit},
				ParentFile: &parent,
				OrderBy:    SortBySize,
				OrderDir:   c.orderDir,
			}, user)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range res.Payload {
				got = append(got, f.Name)
			}
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("size orderDir=%q across pages: got %v, want %v", c.orderDir, got, c.want)
		}
	}
}

// --- type filter coercion ---------------------------------------------------

func TestListFilesSortBySizeReturnsFilesOnly(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	// a directory row large enough in bytes to top a mixed size listing
	insertSortTestSubDir(t, db, dir, "zzz-subdir", user.UserNo, 9999)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)

	want := []string{"large.bin", "small.bin"} // largest first, never a directory
	all := FileTypeFile
	empty := ""
	dirFilter := FileTypeDir
	cases := []struct {
		name     string
		fileType *string
	}{
		{"absent", nil},
		{"all", &empty},
		{"directory", &dirFilter},
		{"file", &all},
	}
	for _, c := range cases {
		payload := listSortedFiles(t, db, dir, ListFileReq{OrderBy: "size", FileType: c.fileType}, user)
		names := make([]string, 0, len(payload))
		for _, f := range payload {
			if f.FileType != FileTypeFile {
				t.Fatalf("fileType=%s: got a %s row in a size listing", c.name, f.FileType)
			}
			names = append(names, f.Name)
		}
		if !reflect.DeepEqual(names, want) {
			t.Fatalf("fileType=%s: got %v, want %v", c.name, names, want)
		}
	}
}

// --- keyword search ---------------------------------------------------------

func TestListFilesSortBySizeRanksKeywordMatches(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "sizezzmarker-small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "sizezzmarker-big.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "unrelated.bin", user.UserNo, 500)

	keyword := "sizezzmarker"
	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", Filename: &keyword}, user)
	want := []string{"sizezzmarker-big.bin", "sizezzmarker-small.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size with keyword: got %v, want %v", got, want)
	}
}

// --- comic directories ------------------------------------------------------

func TestListFilesComicDirSortBySizeReturnsNameAsc(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, true)

	// insertion and PK order are c, a, b, not name order; size order would be
	// b, c, a, and the directory would sort last by name
	insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 900)
	insertSortTestSubDir(t, db, dir, "zzz-subdir", user.UserNo, 9999)

	want := []string{"a.bin", "b.bin", "c.bin"}
	for _, d := range []string{"", "asc", "desc"} {
		payload := listSortedFiles(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: d}, user)
		names := make([]string, 0, len(payload))
		for _, f := range payload {
			if f.FileType != FileTypeFile {
				t.Fatalf("comic dir size orderDir=%q: got a %s row in a files listing", d, f.FileType)
			}
			names = append(names, f.Name)
		}
		if !reflect.DeepEqual(names, want) {
			t.Fatalf("comic dir size orderDir=%q: got %v, want %v", d, names, want)
		}
	}
}

// --- file position ----------------------------------------------------------

// calcSortedPosition asks for the page of fileKey under the given sort and the
// request's raw type filter, the way the position endpoint asks.
func calcSortedPosition(t *testing.T, db *gorm.DB, fileKey, dir, orderBy, orderDir, fileType string, limit int, user flow.User) int {
	t.Helper()
	p, err := CalcFilePosition(miso.EmptyRail(), db, FilePositionReq{
		FileKey:    fileKey,
		ParentFile: dir,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderDir:   orderDir,
		FileType:   fileType,
	}, user)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCalcFilePositionSortBySizeBothDirections(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: small=1, large=2, medium=3; inserted in neither size nor id order
	smallKey := insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	largeKey := insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	mediumKey := insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	// size desc: large, medium, small; size asc: small, medium, large
	cases := []struct {
		orderDir string
		want     map[string]int
	}{
		{"", map[string]int{largeKey: 1, mediumKey: 2, smallKey: 3}},     // omitted direction → desc
		{"desc", map[string]int{largeKey: 1, mediumKey: 2, smallKey: 3}}, // large, medium, small
		{"asc", map[string]int{smallKey: 1, mediumKey: 2, largeKey: 3}},  // small, medium, large
	}
	for _, c := range cases {
		for key, want := range c.want {
			got := calcSortedPosition(t, db, key, dir, SortBySize, c.orderDir, "", 1, user)
			if got != want {
				t.Fatalf("size orderDir=%q: got page %v for %s, want %v", c.orderDir, got, key, want)
			}
		}
	}
}

func TestCalcFilePositionSortBySizeBreaksTiesById(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: c=1, a=2, b=3; all the same size
	cKey := insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	aKey := insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 500)
	bKey := insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 500)

	// size desc falls back to id desc: b, a, c; asc: c, a, b
	cases := []struct {
		orderDir string
		want     map[string]int
	}{
		{"", map[string]int{bKey: 1, aKey: 2, cKey: 3}},
		{"desc", map[string]int{bKey: 1, aKey: 2, cKey: 3}},
		{"asc", map[string]int{cKey: 1, aKey: 2, bKey: 3}},
	}
	for _, c := range cases {
		for key, want := range c.want {
			got := calcSortedPosition(t, db, key, dir, SortBySize, c.orderDir, "", 1, user)
			if got != want {
				t.Fatalf("equal sizes orderDir=%q: got page %v for %s, want %v", c.orderDir, got, key, want)
			}
		}
	}
}

func TestCalcFilePositionSortBySizeCoercesFileType(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	smallKey := insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	largeKey := insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	// a directory row big enough to top a mixed size order; the size sort must
	// never count it
	insertSortTestSubDir(t, db, dir, "zzz-subdir", user.UserNo, 9999)

	// files only, size desc: large (1), small (2) regardless of the asked type
	for _, fileType := range []string{"", FileTypeDir, FileTypeFile} {
		if got, want := calcSortedPosition(t, db, largeKey, dir, SortBySize, "", fileType, 1, user), 1; got != want {
			t.Fatalf("fileType=%q: got page %v for large.bin, want %v", fileType, got, want)
		}
		if got, want := calcSortedPosition(t, db, smallKey, dir, SortBySize, "", fileType, 1, user), 2; got != want {
			t.Fatalf("fileType=%q: got page %v for small.bin, want %v", fileType, got, want)
		}
	}
}

// listingPageOf walks the listing at the given page size and returns the page
// that holds fileKey, failing when the listing never shows it.
func listingPageOf(t *testing.T, db *gorm.DB, dir, orderBy, orderDir string, fileType *string, fileKey string, limit int, user flow.User) int {
	t.Helper()
	parent := dir
	for page := 1; page < 100; page++ {
		res, err := ListFiles(miso.EmptyRail(), db, ListFileReq{
			Page:       miso.Paging{Page: page, Limit: limit},
			ParentFile: &parent,
			OrderBy:    orderBy,
			OrderDir:   orderDir,
			FileType:   fileType,
		}, user)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Payload) == 0 {
			break
		}
		for _, f := range res.Payload {
			if f.Uuid == fileKey {
				return page
			}
		}
	}
	t.Fatalf("file %s never appears in the listing", fileKey)
	return 0
}

func TestCalcFilePositionSortBySizeMatchesListingPage(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	smallKey := insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	largeKey := insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	mediumKey := insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	empty := ""
	dirFilter := FileTypeDir
	cases := []struct {
		orderDir string
		fileType *string
	}{
		{"", nil},        // omitted direction, type filter absent
		{"asc", nil},     // smallest first
		{"desc", &empty}, // "all" type filter coerced to files
		{"", &dirFilter}, // directory type filter coerced to files
	}
	limit := 2
	for _, c := range cases {
		fileType := ""
		if c.fileType != nil {
			fileType = *c.fileType
		}
		for _, key := range []string{smallKey, largeKey, mediumKey} {
			want := listingPageOf(t, db, dir, SortBySize, c.orderDir, c.fileType, key, limit, user)
			got := calcSortedPosition(t, db, key, dir, SortBySize, c.orderDir, fileType, limit, user)
			if got != want {
				t.Fatalf("size orderDir=%q fileType=%q: %s position is %v, but the listing shows it on page %v", c.orderDir, fileType, key, got, want)
			}
		}
	}
}

func TestCalcFilePositionComicDirSortBySizeUsesNameAsc(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, true)

	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 900)
	cKey := insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	// a subdirectory that sorts before c.bin by name, so a non-coerced count
	// would push c.bin one page further than the listing does
	insertSortTestSubDir(t, db, dir, "aa-subdir", user.UserNo, 9999)

	// comic reading order is name asc, files only: a, b, c
	for _, d := range []string{"", "asc", "desc"} {
		got := calcSortedPosition(t, db, cKey, dir, SortBySize, d, "", 1, user)
		if want := listingPageOf(t, db, dir, SortBySize, d, nil, cKey, 1, user); got != want {
			t.Fatalf("comic size orderDir=%q: position %v, listing page %v", d, got, want)
		}
		if got != 3 {
			t.Fatalf("comic size orderDir=%q: got page %v for c.bin, want 3 (name asc)", d, got)
		}
	}
}

// A size position is resolved from the rows themselves: the calculation neither
// records nor reads a resume cursor.
func TestCalcFilePositionSortBySizeRecordsNoResumeCursor(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	smallKey := insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)

	if got, want := calcSortedPosition(t, db, smallKey, dir, SortBySize, "", "", 1, user), 2; got != want {
		t.Fatalf("got page %v for small.bin, want %v", got, want)
	}

	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM dir_resume_cursor WHERE dir_key = ?", dir).Scan(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("size position recorded %v resume cursor row(s), want 0", count)
	}
}
