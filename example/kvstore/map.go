package main

import (
	"solod.dev/so/c"
	"solod.dev/so/errors"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/strings"
	"solod.dev/sqlite/libsqlite3"
)

var (
	ErrCreate   = errors.New("kvstore: create schema failed")
	ErrExec     = errors.New("kvstore: exec failed")
	ErrNotFound = errors.New("kvstore: not found")
	ErrPrepare  = errors.New("kvstore: prepare failed")
)

const (
	sqlCreate = "create table if not exists kv (key text primary key, val)"
	sqlGet    = "select val from kv where key = ?"
	sqlSet    = "insert or replace into kv (key, val) values (?, ?)"
	sqlDelete = "delete from kv where key = ?"
)

// Store is a simple key-value store backed by an SQLite database.
type Store struct {
	db *libsqlite3.Sqlite3
}

// NewStore creates a new Store using the provided connection string.
// It opens a connection to the SQLite database and creates the underlying
// key-value table if it does not already exist.
//
// Call [Store.Close] when done to release resources.
func NewStore(connStr string) (Store, error) {
	var db *libsqlite3.Sqlite3
	rc := libsqlite3.Open(connStr, &db)
	if rc != libsqlite3.SQLITE_OK {
		return Store{}, ErrCreate
	}

	rc = libsqlite3.Exec(db, sqlCreate, nil, nil, nil)
	if rc != libsqlite3.SQLITE_OK {
		libsqlite3.Close(db)
		return Store{}, ErrCreate
	}
	return Store{db}, nil
}

// GetInt returns the integer value associated with the specified key.
func (s *Store) GetInt(key string) (int, error) {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlGet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return 0, ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	rc = libsqlite3.Step(stmt)
	if rc == libsqlite3.SQLITE_DONE {
		return 0, ErrNotFound
	}
	if rc != libsqlite3.SQLITE_ROW {
		return 0, ErrExec
	}

	result := int(libsqlite3.Column_int64(stmt, 0))
	return result, nil
}

// GetFloat64 returns the float64 value associated with the specified key.
func (s *Store) GetFloat64(key string) (float64, error) {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlGet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return 0, ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	rc = libsqlite3.Step(stmt)
	if rc == libsqlite3.SQLITE_DONE {
		return 0, ErrNotFound
	}
	if rc != libsqlite3.SQLITE_ROW {
		return 0, ErrExec
	}

	result := libsqlite3.Column_double(stmt, 0)
	return result, nil
}

// GetString returns the string value associated with the specified key.
// The returned string is allocated; the caller owns it.
func (s *Store) GetString(alloc mem.Allocator, key string) (string, error) {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlGet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return "", ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	rc = libsqlite3.Step(stmt)
	if rc == libsqlite3.SQLITE_DONE {
		return "", ErrNotFound
	}
	if rc != libsqlite3.SQLITE_ROW {
		return "", ErrExec
	}

	text := c.PtrAs[c.ConstChar](libsqlite3.Column_text(stmt, 0))
	src := c.String(text)
	result := strings.Clone(alloc, src)
	return result, nil
}

// GetByte returns the raw blob value associated with the specified key.
// The returned slice is allocated; the caller owns it.
func (s *Store) GetByte(alloc mem.Allocator, key string) ([]byte, error) {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlGet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return []byte{}, ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	rc = libsqlite3.Step(stmt)
	if rc == libsqlite3.SQLITE_DONE {
		return []byte{}, ErrNotFound
	}
	if rc != libsqlite3.SQLITE_ROW {
		return []byte{}, ErrExec
	}

	ptr := libsqlite3.Column_blob(stmt, 0).(*byte)
	n := libsqlite3.Column_bytes(stmt, 0)
	src := c.Bytes(ptr, int(n))
	result := slices.Clone(alloc, src)
	return result, nil
}

// SetInt stores an integer value for the specified key.
func (s *Store) SetInt(key string, val int) error {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlSet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	libsqlite3.Bind_int64(stmt, 2, c.LongLong(val))

	rc = libsqlite3.Step(stmt)
	if rc != libsqlite3.SQLITE_DONE {
		return ErrExec
	}
	return nil
}

// SetFloat64 stores a float64 value for the specified key.
func (s *Store) SetFloat64(key string, val float64) error {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlSet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	libsqlite3.Bind_double(stmt, 2, val)

	rc = libsqlite3.Step(stmt)
	if rc != libsqlite3.SQLITE_DONE {
		return ErrExec
	}
	return nil
}

// SetString stores a string value for the specified key.
func (s *Store) SetString(key string, val string) error {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlSet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	libsqlite3.Bind_text(stmt, 2, val, c.Int(len(val)), nil)

	rc = libsqlite3.Step(stmt)
	if rc != libsqlite3.SQLITE_DONE {
		return ErrExec
	}
	return nil
}

// SetByte stores a blob value for the specified key.
func (s *Store) SetByte(key string, val []byte) error {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlSet, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	libsqlite3.Bind_blob(stmt, 2, val, c.Int(len(val)), nil)

	rc = libsqlite3.Step(stmt)
	if rc != libsqlite3.SQLITE_DONE {
		return ErrExec
	}
	return nil
}

// Delete removes the entry with the specified key.
func (s *Store) Delete(key string) error {
	var stmt *libsqlite3.Stmt
	rc := libsqlite3.Prepare_v2(s.db, sqlDelete, -1, &stmt, nil)
	if rc != libsqlite3.SQLITE_OK {
		return ErrPrepare
	}
	defer libsqlite3.Finalize(stmt)

	libsqlite3.Bind_text(stmt, 1, key, c.Int(len(key)), nil)
	rc = libsqlite3.Step(stmt)

	if rc != libsqlite3.SQLITE_DONE {
		return ErrExec
	}
	return nil
}

// Close releases resources associated with the Store.
func (s *Store) Close() {
	libsqlite3.Close(s.db)
}
