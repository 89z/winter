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
   db, err := sql.Open("sqlite3", file)
   if err != nil {
      return Tx{}, err
   }
   tx, err := db.Begin()
   if err != nil {
      return Tx{}, err
   }
   return Tx{tx}, nil
}

func (tx Tx) Delete(query string, arg ...interface{}) error {
   fmt.Println(arg)
   _, err := tx.Exec("delete from " + query, arg...)
   return err
}

func (tx Tx) Insert(query string, arg ...interface{}) (int64, error) {
   fmt.Println(arg)
   res, err := tx.Exec("insert into " + query, arg...)
   if err != nil {
      return 0, err
   }
   return res.LastInsertId()
}

func (tx Tx) Update(query string, arg ...interface{}) error {
   fmt.Println(arg)
   _, err := tx.Exec("update " + query, arg...)
   return err
}
