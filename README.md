# sqlite

Solod bindings for [SQLite](https://sqlite.org), a small, fast, and self-contained SQL database engine.

## Usage

1. Install SQLite for your operating system.

2. Install the Solod bindings.

```
go get solod.dev/sqlite@latest
```

3. Use it in your code.

```go
package main

import (
	"solod.dev/so/c"
	"solod.dev/so/fmt"
	"solod.dev/so/os"
	"solod.dev/sqlite/libsqlite3"
)

func main() {
	var db *libsqlite3.Sqlite3
	if rc := libsqlite3.Open(":memory:", &db); rc != libsqlite3.SQLITE_OK {
		fmt.Println("Can't open database:", c.String(libsqlite3.Errmsg(db)))
		os.Exit(1)
	}
	defer libsqlite3.Close(db)
	fmt.Println("SQLite version:", c.String(libsqlite3.Libversion()))
}
```

## Examples

[Simple query](example/simple/main.go)

[Key-value store](example/kvstore/main.go)
