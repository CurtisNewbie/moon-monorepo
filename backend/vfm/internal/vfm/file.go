package vfm

import (
	"os"

	"github.com/curtisnewbie/miso/miso"
)

func MakeTempDirs(rail miso.Rail) error {
	dir := miso.GetPropStr(PropTempPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	rail.Infof("MkdirAll %v finished", dir)
	return nil
}

func TempFilePath(tempTkn string) string {
	return miso.GetPropStr(PropTempPath) + "/" + tempTkn
}
