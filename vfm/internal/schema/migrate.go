package schema

import (
	"embed"

	"github.com/curtisnewbie/miso/middleware/svc"
)

//go:embed scripts/*.sql
var SchemaFs embed.FS

const (
	BaseDir     = "scripts"
	StartingVer = "v0.1.17.sql"
)

func init() {
	svc.ExcludeSchemaFile("schema.sql")
}

// starting from v0.1.18, let svc manages the schema migration
func EnableSchemaMigrate() {
	svc.EnableSchemaMigrate(SchemaFs, BaseDir, StartingVer)
}
