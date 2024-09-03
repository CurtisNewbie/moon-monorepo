package vfm

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TagA        = "a" // bookmark file <A>
	AttrHref    = "href"
	AttrAddDate = "add_date"
	AttrIcon    = "icon"
)

type NetscapeBookmarkFile struct {
	Bookmarks []Bookmark
}

func (n *NetscapeBookmarkFile) Add(b Bookmark) {
	n.Bookmarks = append(n.Bookmarks, b)
}

type Bookmark struct {
	Name    string
	Href    string
	Icon    string
	AddDate string
}

func (b Bookmark) String() string {
	return fmt.Sprintf("Bookmark{\n\tName: %v,\n\tHref: %v,\n\tIcon: %v,\n\tAddDate: %v\n}", b.Name, b.Href, b.Icon, b.AddDate)
}

func ParseNetscapeBookmark(rail miso.Rail, body io.Reader) (NetscapeBookmarkFile, error) {
	bookmarkFile := NetscapeBookmarkFile{Bookmarks: []Bookmark{}}

	z := html.NewTokenizer(body)
	curr := Bookmark{}

	for {
		ttype := z.Next()
		bname, isAttr := z.TagName()
		name := string(bname)

		textB := z.Text()
		text := string(textB)

		attr := map[string]string{}

		for {
			attrKeyB, attrValB, more := z.TagAttr()
			k := string(attrKeyB)
			v := string(attrValB)
			attr[k] = v
			if !more {
				break
			}
		}

		// rail.Debugf("tokenType: %v, text: %v, name: %v, isAttr: %v, attr: %v",
		// 	ttype, text, name, isAttr, attr)

		switch ttype {
		case html.ErrorToken:
			err := z.Err()
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return bookmarkFile, err
		case html.TextToken:
			if !isAttr {
				curr.Name = text
			}
		case html.StartTagToken:
			if name == TagA && isAttr {
				if v, ok := attr[AttrHref]; ok {
					curr.Href = v
				}
				if v, ok := attr[AttrAddDate]; ok {
					curr.AddDate = v
				}
				if v, ok := attr[AttrIcon]; ok {
					curr.Icon = v
				}
			}
		case html.EndTagToken:
			if name == TagA {
				bookmarkFile.Add(curr)
				curr = Bookmark{}
			}
		}
	}
}

func TransferTmpFile(rail miso.Rail, reader io.Reader) (string, error) {
	path := TempFilePath(util.RandAlpha(15))

	f, err := os.Create(path)
	if err != nil {
		return "", ErrUploadFiled.WithInternalMsg("create file failed, path: %v, %v", path, err)
	}
	if _, err := io.Copy(f, reader); err != nil {
		return "", ErrUploadFiled.WithInternalMsg("transfer file failed, path: %v, %v", path, err)
	}
	rail.Infof("Transferred file to path: %v", path)
	return path, nil
}

func ProcessUploadedBookmarkFile(rail miso.Rail, path string, user common.User) error {
	rail.Infof("User '%v' parse bookmark file, tmpFile: %v", user.Username, path)
	file, err := os.Open(path)
	if err != nil {
		return ErrUnknown.WithInternalMsg("open temp file failed, path: %v", path)
	}

	bookmarkFile, err := ParseNetscapeBookmark(rail, file)
	if err != nil {
		return ErrUnknown.WithInternalMsg("open temp file failed, path: %v", path)
	}

	go func() {
		rail := rail.NextSpan()
		err := SaveBookmarks(rail, mysql.GetMySQL(), bookmarkFile, user)
		if err != nil {
			rail.Errorf("failed to save bookmark, user: %s, %v", user.Username, err)
		}
	}()
	return nil
}

type NewBookmark struct {
	UserNo string
	Icon   string
	Name   string
	Href   string
	Md5    string
}

func SaveBookmarks(rail miso.Rail, tx *gorm.DB, bookmarkFile NetscapeBookmarkFile, user common.User) error {

	bookmarks := make([]NewBookmark, 0, len(bookmarkFile.Bookmarks))
	for i := range bookmarkFile.Bookmarks {
		bm := bookmarkFile.Bookmarks[i]
		md5 := BookmarkMd5(bm)

		var id int
		t := tx.Raw(`SELECT id FROM bookmark_blacklist WHERE user_no = ? and md5 = ?`, user.UserNo, md5).Scan(&id)
		if t.RowsAffected > 0 {
			rail.Infof("bookmark in blacklist, ignored, userNo: %s, md5: %s, name: %s", user.UserNo, md5, bm.Name)
			continue
		}

		bookmarks = append(bookmarks, NewBookmark{
			UserNo: user.UserNo,
			Icon:   bm.Icon,
			Name:   bm.Name,
			Href:   bm.Href,
			Md5:    md5,
		})
	}
	// rail.Debugf("bookmarks: %+v", bookmarks)

	err := tx.Table("bookmark").Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(bookmarks, 100).Error
	if err != nil {
		return fmt.Errorf("failed to insert bookmark, %v", err)
	}

	return nil
}

func BookmarkMd5(bm Bookmark) string {
	s := fmt.Sprintf("NA%vHR%vIC%v", bm.Name, bm.Href, bm.Icon)
	chksum := md5.Sum([]byte(s))
	return hex.EncodeToString(chksum[:])
}

type ListedBookmark struct {
	Id     int64
	UserNo string
	Name   string
	Href   string
	Icon   string
}

func ListBookmarks(rail miso.Rail, tx *gorm.DB, req ListBookmarksReq, userNo string) (any, error) {
	return mysql.NewPageQuery[ListedBookmark]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			if req.Blacklisted {
				tx = tx.Table("bookmark_blacklist")
			} else {
				tx = tx.Table("bookmark")
			}
			tx = tx.Where("user_no = ?", userNo)
			if req.Name != nil && *req.Name != "" {
				tx = tx.Where("name like ?", "%"+*req.Name+"%")
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id, user_no, name, href, icon").
				Order("id DESC").
				Offset(req.Paging.GetOffset()).
				Limit(req.Paging.GetLimit())
		}).
		Exec(rail, tx)
}

type RemoveBookmarkInf struct {
	UserNo string
	Icon   string
	Name   string
	Href   string
	Md5    string
}

func RemoveBookmark(rail miso.Rail, tx *gorm.DB, id int64, userNo string) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var b RemoveBookmarkInf
		tx = tx.Raw(`SELECT * FROM bookmark WHERE id = ?`, id).Scan(&b)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected < 1 {
			return miso.NewErrf("Bookmark not found")
		}
		if b.UserNo != userNo {
			return miso.ErrNotPermitted
		}

		err := tx.Exec(`INSERT IGNORE INTO bookmark_blacklist (user_no, icon, name, href, md5) VALUES (?,?,?,?,?)`, b.UserNo, b.Icon, b.Name, b.Href, b.Md5).Error
		if err != nil {
			return err
		}
		return tx.Exec("DELETE FROM bookmark WHERE user_no = ? AND id = ?", userNo, id).Error
	})
}

func RemoveBookmarkBlacklist(rail miso.Rail, tx *gorm.DB, id int64, userNo string) error {
	return tx.Exec("DELETE FROM bookmark_blacklist WHERE user_no = ? AND id = ?", userNo, id).Error
}
