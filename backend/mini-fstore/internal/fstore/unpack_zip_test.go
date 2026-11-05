package fstore

import (
	"strings"
	"testing"

	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/miso/miso"
)

func TestUnpackZipSkipsDirectoryEntries(t *testing.T) {
	miso.SetProp(config.PropStorageDir, t.TempDir())

	fileId := "file_zip_dir_entries"
	writeTestZipAt(t, GenStoragePath(fileId), []testZipEntry{
		{Name: "Documents/"},
		{Name: "Documents/doc/"},
		{Name: "Documents/doc/a.txt", Body: []byte("a")},
		{Name: "Documents/b.txt", Body: []byte("b")},
	})

	entries, err := UnpackZip(miso.EmptyRail(), File{FileId: fileId}, t.TempDir())
	if err != nil {
		t.Fatalf("UnpackZip returned error, %v", err)
	}

	want := []string{"Documents/doc/a.txt", "Documents/b.txt"}
	if len(entries) != len(want) {
		t.Fatalf("UnpackZip returned %d entries, want %d: %+v", len(entries), len(want), entries)
	}
	for i, e := range entries {
		if e.Name != want[i] {
			t.Errorf("entry[%d].Name = %q, want %q", i, e.Name, want[i])
		}
		if strings.HasSuffix(e.Name, "/") {
			t.Errorf("entry[%d].Name = %q, want no directory entry", i, e.Name)
		}
	}
}
