package snow

import (
   "database/sql"
   "fmt"
)

func Update(o *sql.DB, query string, args ...interface{}) error {
   fmt.Println(args)
   _, e := o.Exec("update " + query, args...)
   return e
}

func Delete(o *sql.DB, query string, args ...interface{}) error {
   fmt.Println(args)
   _, e := o.Exec("delete from " + query, args...)
   return e
}

func Insert(open_o *sql.DB, query string, args ...interface{}) (int64, error) {
   fmt.Println(args)
   exec_o, e := open_o.Exec("insert into " + query, args...)
   if e != nil {
      return 0, e
   }
   return exec_o.LastInsertId()
}
