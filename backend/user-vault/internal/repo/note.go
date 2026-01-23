package repo

import (
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/idutil"
	"gorm.io/gorm"
)

const (
	TableNote = "note"
)

var (
	ErrNoteNotFound = errs.NewErrfCode("NOTE_NOT_FOUND", "Note not found")
)

type SaveNoteReq struct {
	Title   string `valid:"trim,notEmpty" json:"title"`
	Content string `json:"content"`
}

func SaveNote(rail miso.Rail, db *gorm.DB, snr SaveNoteReq, user flow.User) error {
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
	RecordId string `valid:"notEmpty" json:"recordId"`
	Title    string `valid:"trim,notEmpty" json:"title"`
	Content  string `json:"content"`
}

func UpdateNote(rail miso.Rail, db *gorm.DB, unr UpdateNoteReq, user flow.User) error {
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

func DeleteNote(rail miso.Rail, db *gorm.DB, recordId string, user flow.User) error {
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
	RecordId  string    `json:"recordId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserNo    string    `json:"userNo"`
	CreatedAt atom.Time `json:"createdAt"`
	UpdatedAt atom.Time `json:"updatedAt"`
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
	Keywords string      `json:"keywords"`
	Paging   miso.Paging `json:"paging"`
}

func ListNotes(rail miso.Rail, db *gorm.DB, req ListNoteReq, user flow.User) (miso.PageRes[Note], error) {
	return dbquery.NewPagedQuery[Note](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table(TableNote).
				Eq("user_no", user.UserNo).
				Eq("deleted", false).
				WhereIf(req.Keywords != "", "MATCH (title, content) AGAINST (? IN BOOLEAN MODE)", req.Keywords)
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.SelectCols(Note{})
		}).
		Scan(rail, req.Paging)
}
