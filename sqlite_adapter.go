package reactiverecord

import (
  "fmt"
  "github.com/eaigner/jet"
  _ "github.com/mattn/go-sqlite3"
  // "strconv"
  // "strings"
  // "time"
)

type SQLiteAdapter struct {
  db *jet.Db
}

func NewSqliteAdapter(params ...string) (*SQLiteAdapter, error) {
  var err error
  databaseName := "data"
  store := new(SQLiteAdapter)
  if len(params) > 0 {
    databaseName = params[0]
  }
  store.db, err = jet.Open("sqlite3", databaseName)
  if err != nil {
    fmt.Printf("Error connecting to database %s\n", err)
  }
  return store, err
}

func (a *SQLiteAdapter) Query(query string) jet.Runnable {
  return a.db.Query(query)
}
