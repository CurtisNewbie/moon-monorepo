package vfm

import (
	"testing"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
)

func TestListVerFile(t *testing.T) {
	corePreTest(t)
	rail := miso.EmptyRail()
	files, err := ListVerFile(rail, mysql.GetMySQL(), ApiListVerFileReq{}, testUser())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", files)
}

func TestDelVerFile(t *testing.T) {
	corePreTest(t)
	rail := miso.EmptyRail()
	verFileId := "verf_1224911895724032158144"

	err := DelVerFile(rail, mysql.GetMySQL(),
		ApiDelVerFileReq{VerFileId: verFileId},
		testUser())

	if err != nil {
		t.Fatal(err)
	}
}
