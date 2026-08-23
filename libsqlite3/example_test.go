package libsqlite3_test

import (
	"solod.dev/so/c"
	"solod.dev/so/fmt"
	"solod.dev/sqlite/libsqlite3"
)

func ExampleOpen() {
	var db *libsqlite3.Sqlite3
	if rc := libsqlite3.Open(":memory:", &db); rc != libsqlite3.SQLITE_OK {
		fmt.Println("Can't open database:", c.String(libsqlite3.Errmsg(db)))
		return
	}
	defer libsqlite3.Close(db)
	fmt.Println("SQLite version:", c.String(libsqlite3.Libversion()))
}
