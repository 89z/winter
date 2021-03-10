package winter

import (
   "database/sql"
   "fmt"
   _ "github.com/mattn/go-sqlite3"
)

type Rows = sql.Rows

type Tx struct {
   *sql.Tx
}

func NewTx(file string) (Tx, error) {
   db, e := sql.Open("sqlite3", file)
   if e != nil {
      return Tx{}, e
   }
   tx, e := db.Begin()
   if e != nil {
      return Tx{}, e
   }
   return Tx{tx}, nil
}

func (tx Tx) Delete(query string, arg ...interface{}) error {
   fmt.Println(arg)
   _, e := tx.Exec("delete from " + query, arg...)
   return e
}

func (tx Tx) Insert(query string, arg ...interface{}) (int64, error) {
   fmt.Println(arg)
   res, e := tx.Exec("insert into " + query, arg...)
   if e != nil {
      return 0, e
   }
   return res.LastInsertId()
}

func (tx Tx) Update(query string, arg ...interface{}) error {
   fmt.Println(arg)
   _, e := tx.Exec("update " + query, arg...)
   return e
}
