package main

import (
   "database/sql"
   "fmt"
   "time"
)

func Exec(open_o *sql.DB, query string, args ...interface{}) (int64, error) {
   fmt.Println(args)
   exec_o, e := open_o.Exec(query, args...)
   if e != nil {
      return 0, e
   }
   return exec_o.LastInsertId()
}

func Note(f float64) string {
   d := time.Duration(f) * time.Millisecond
   if d < 179_500 * time.Millisecond {
      return "short"
   }
   if d > 15 * time.Minute {
      return "long"
   }
   return ""
}
