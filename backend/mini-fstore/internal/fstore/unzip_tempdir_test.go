package fstore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/miso/miso"
)

// UnpackZip writes every entry into the temp dir, so the dir has to be
// traversable by its owner. Creating it with a file mode (0666) yields a dir
// without the execute bit and the first entry fails with a permission error.
func TestNewUnzipTempDirIsWritable(t *testing.T) {
	miso.SetProp(config.PropTempDir, t.TempDir())

	dir, err := newUnzipTempDir("file_tempdir_mode")
	if err != nil {
		t.Fatalf("newUnzipTempDir returned error, %v", err)
	}
	defer os.RemoveAll(dir)

	f, err := os.Create(filepath.Join(dir, "ZIPENTRY"))
	if err != nil {
		t.Fatalf("temp dir %v is not writable, %v", dir, err)
	}
	defer f.Close()
}
