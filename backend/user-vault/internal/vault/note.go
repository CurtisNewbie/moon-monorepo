package vault

import (
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/idutil"
	"gorm.io/gorm"
)

const (
	TableNote = "note"
)

type SaveNoteReq struct {
	Title   string
	Content string
}

func DBSaveNote(rail miso.Rail, db *gorm.DB, snr SaveNoteReq, user common.User) error {
	return dbquery.NewQuery(rail, db).
		Table(TableNote).
		CreateAny(struct {
			RecordId  string
			Title     string
			Content   string
			UserNo    string
			CreatedBy string
		}{
			idutil.Id("note_"),
			snr.Title,
			snr.Content,
			user.UserNo,
			user.Username,
		})
}

type UpdateNoteReq struct {
	RecordId string
	Title    string
	Content  string
}

func DBUpdateNote(rail miso.Rail, db *gorm.DB, unr UpdateNoteReq, user common.User) error {
	return dbquery.NewQuery(rail, db).
		Table(TableNote).
		SetCols(struct {
			Title     string
			Content   string
			UpdatedBy string
		}{
			unr.Title,
			unr.Content,
			user.Username,
		}).
		Eq("record_id", unr.RecordId).
		Eq("user_no", user.UserNo).
		UpdateAny()
}

func DBUDeleteNote(rail miso.Rail, db *gorm.DB, recordId string, user common.User) error {
	return dbquery.NewQuery(rail, db).
		Table(TableNote).
		SetCols(struct {
			Deleted   bool
			UpdatedBy string
		}{
			true,
			user.Username,
		}).
		Eq("record_id", recordId).
		Eq("user_no", user.UserNo).
		Eq("deleted", false).
		UpdateAny()
}
