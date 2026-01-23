package note

import (
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/repo"
	"gorm.io/gorm"
)

func NewNoteLock(rail miso.Rail, recordId string) *redis.RLock {
	return redis.NewRLockf(rail, "user-vault:note:record-id:%v", recordId)
}

func ListNotes(rail miso.Rail, db *gorm.DB, req repo.ListNoteReq, user flow.User) (miso.PageRes[repo.Note], error) {
	return repo.ListNotes(rail, db, req, user)
}

func UpdateNote(rail miso.Rail, db *gorm.DB, req repo.UpdateNoteReq, user flow.User) error {
	lock := NewNoteLock(rail, req.RecordId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	_, err := repo.FindNote(rail, db, req.RecordId, user.UserNo)
	if err != nil {
		return err
	}
	return repo.UpdateNote(rail, db, req, user)
}

func DeleteNote(rail miso.Rail, db *gorm.DB, recordId string, user flow.User) error {
	lock := NewNoteLock(rail, recordId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	_, err := repo.FindNote(rail, db, recordId, user.UserNo)
	if err != nil {
		return err
	}
	return repo.DeleteNote(rail, db, recordId, user)
}

func SaveNote(rail miso.Rail, db *gorm.DB, snr repo.SaveNoteReq, user flow.User) error {
	return repo.SaveNote(rail, db, snr, user)
}
