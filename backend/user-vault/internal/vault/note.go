package vault

import (
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/idutil"
	"gorm.io/gorm"
)

const (
	TableNote = "note"
)

var (
	ErrNoteNotFound = errs.NewErrfCode("NOTE_NOT_FOUND", "Note not found")
)

func NewNoteLock(rail miso.Rail, recordId string) *redis.RLock {
	return redis.NewRLockf(rail, "user-vault:note:record-id:%v", recordId)
}

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

func DBDeleteNote(rail miso.Rail, db *gorm.DB, recordId string, user common.User) error {
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

type Note struct {
	RecordId string
	Title    string
	Content  string
	UserNo   string
}

func FindNote(rail miso.Rail, db *gorm.DB, recordId string, userNo string) (Note, error) {
	var n Note
	ok, err := dbquery.NewQuery(rail, db).
		Table(TableNote).
		SelectCols(n).
		Eq("record_id", recordId).
		Eq("user_no", userNo).
		Eq("deleted", false).
		ScanAny(&n)
	if err != nil {
		return n, err
	}
	if !ok {
		return n, ErrNoteNotFound.New()
	}
	return n, nil
}

type ListNoteReq struct {
	Keywords string
	Paging   miso.Paging
}

func ListNotes(rail miso.Rail, db *gorm.DB, req ListNoteReq, user common.User) (miso.PageRes[Note], error) {
	return dbquery.NewPagedQuery[Note](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table(TableNote).
				Eq("user_no", user.UserNo).
				Eq("deleted", false).
				WhereIf(req.Keywords != "", "title LIKE ? OR content LIKE ?", "%"+req.Keywords+"%") // TODO: replace like with match()
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.SelectCols(Note{})
		}).
		Scan(rail, req.Paging)
}

func UpdateNote(rail miso.Rail, db *gorm.DB, req UpdateNoteReq, user common.User) error {
	lock := NewNoteLock(rail, req.RecordId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	_, err := FindNote(rail, db, req.RecordId, user.UserNo)
	if err != nil {
		return err
	}
	return DBUpdateNote(rail, db, req, user)
}

func DeleteNote(rail miso.Rail, db *gorm.DB, recordId string, user common.User) error {
	lock := NewNoteLock(rail, recordId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	_, err := FindNote(rail, db, recordId, user.UserNo)
	if err != nil {
		return err
	}
	return DBDeleteNote(rail, db, recordId, user)
}
