package vfm

import (
	"os"
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
)

func TestNetscapeBookmark(t *testing.T) {
	testFile := "bookmarks_2023_10_13.html"
	f, err := os.Open(testFile)
	if err != nil {
		t.Skipf("skipping: bookmark fixture %q is not available, %v", testFile, err)
	}
	defer f.Close()

	rail := miso.EmptyRail()
	// rail.SetLogLevel("debug")
	nb, err := ParseNetscapeBookmark(rail, f)
	if err != nil {
		t.Logf("failed to parse netscape bookmark, %v, typ: %v", err, reflect.TypeOf(err))
		t.FailNow()
	}

	for i := range nb.Bookmarks {
		b := nb.Bookmarks[i]
		rail.Infof("[%v], %v", i, b)
	}
	// t.Logf("bookmark: %+v", nb)
}

func TestProcessUploadedBookmarkFile(t *testing.T) {
	testFile := "bookmarks_2023_10_13.html"
	if _, err := os.Stat(testFile); err != nil {
		t.Skipf("skipping: bookmark fixture %q is not available, %v", testFile, err)
	}

	corePreTest(t)
	rail := miso.EmptyRail()
	miso.SetLogLevel("debug")

	user := flow.User{
		UserNo:   "UE202205142310076187414",
		Username: "zhuangyongj",
	}
	if err := ProcessUploadedBookmarkFile(rail, testFile, user); err != nil {
		t.Log(err)
		t.FailNow()
	}
}
