// A simple key-value store backed by an SQLite database.
//
// Usage:
//
//	make example name=kvstore
//	./build/kvstore -db build/kvstore.db -op set -key name -val alice
//	./build/kvstore -db build/kvstore.db -op get -key name
//	./build/kvstore -db build/kvstore.db -op del -key name
package main

import (
	"solod.dev/so/flag"
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/os"
)

var (
	dbFlag  string
	opFlag  string
	keyFlag string
	valFlag string
)

func main() {
	parseFlags()

	store, err := NewStore(dbFlag)
	check(err)
	defer store.Close()

	switch opFlag {
	case "set":
		err = store.SetString(keyFlag, valFlag)
		check(err)
	case "get":
		val, err := store.GetString(mem.System, keyFlag)
		check(err)
		if err == ErrNotFound {
			fmt.Println("(none)")
		} else {
			fmt.Println(val)
		}
		mem.FreeString(mem.System, val)
	case "del":
		err = store.Delete(keyFlag)
		check(err)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func parseFlags() {
	flag.StringVar(&dbFlag, "db", "kvstore.db", "database file")
	flag.StringVar(&opFlag, "op", "", "operation: get, set, or del")
	flag.StringVar(&keyFlag, "key", "", "key name")
	flag.StringVar(&valFlag, "val", "", "value (for set operation)")
	flag.Parse()
}

func check(err error) {
	if err != nil && err != ErrNotFound {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
