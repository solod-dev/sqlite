// Using the sqlite3 library to execute a SQL statement
// on a database and print the results.
//
// Usage:
//
//	make example name=simple
//	./build/simple build/simple.db "create table data(message text)"
//	./build/simple build/simple.db "insert into data values('Hello, World!')"
//	./build/simple build/simple.db "select * from data"
//
// Source: https://sqlite.org/quickstart.html
package main

import (
	"solod.dev/so/c"
	"solod.dev/so/fmt"
	"solod.dev/so/os"
	"solod.dev/sqlite/libsqlite3"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <database> <sql-statement>\n", os.Args[0])
		os.Exit(1)
	}

	// Open the database.
	var db *libsqlite3.Sqlite3
	rc := libsqlite3.Open(os.Args[1], &db)
	if rc != 0 {
		msg := libsqlite3.Errmsg(db)
		fmt.Fprintf(os.Stderr, "Can't open database: %s\n", c.String(msg))
		os.Exit(1)
	}
	defer libsqlite3.Close(db)

	// Execute the SQL statement.
	var errMsg *c.Char
	rc = libsqlite3.Exec(db, os.Args[2], callback, nil, &errMsg)
	if rc != libsqlite3.SQLITE_OK {
		fmt.Fprintf(os.Stderr, "SQL error: %s\n", c.String(errMsg))
		libsqlite3.Free(errMsg)
	}
}

// callback is called for each row in the result set.
func callback(unused any, nCols c.Int, vals **c.Char, cols **c.Char) c.Int {
	_ = unused
	for i := 0; i < int(nCols); i++ {
		col := *c.PtrAdd(cols, i)
		val := *c.PtrAdd(vals, i)
		if val != nil {
			fmt.Printf("%s = %s\n", c.String(col), c.String(val))
		} else {
			fmt.Printf("%s = NULL\n", c.String(col))
		}
	}
	return 0
}
