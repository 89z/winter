package main

import (
   "database/sql"
   "log"
   _ "github.com/mattn/go-sqlite3"
)

func Insert() {
   open_o, e := sql.Open("sqlite3", "winter.db")
   if e != nil {
      log.Fatal(e)
   }
   s := `
   insert into artist_t (artist_s) values ('April')
   `
   exec_o, e := open_o.Exec(s)
   if e != nil {
      log.Fatal(exec_o, e)
   }
   n, e := exec_o.LastInsertId()
   if e != nil {
      log.Fatal(e)
   }
   println(n)
}
