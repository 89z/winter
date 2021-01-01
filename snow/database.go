package snow

import (
   "database/sql"
   "fmt"
)

func Update(tx sql.Tx, query string, args ...interface{}) error {
   fmt.Println(args)
   _, e := tx.Exec("update " + query, args...)
   return e
}

func Delete(tx sql.Tx, query string, args ...interface{}) error {
   fmt.Println(args)
   _, e := tx.Exec("delete from " + query, args...)
   return e
}

func Insert(tx sql.Tx, query string, args ...interface{}) (int64, error) {
   fmt.Println(args)
   res, e := tx.Exec("insert into " + query, args...)
   if e != nil {
      return 0, e
   }
   return res.LastInsertId()
}
