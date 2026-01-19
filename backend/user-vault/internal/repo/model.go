package repo

import "github.com/curtisnewbie/miso/middleware/dbquery"

func init() {
	dbquery.PrepareCreateModelHook()
	dbquery.PrepareUpdateModelHook()
}
