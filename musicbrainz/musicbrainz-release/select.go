package main

import (
   "database/sql"
   "log"
   _ "github.com/mattn/go-sqlite3"
)

func Select() {
   open_o, e := sql.Open("sqlite3", "winter.db")
   if e != nil {
      log.Fatal(e)
   }
   s := `
   select artist_n
   from artist_t
   where artist_s = 'Cocteau Twins'
   `
   query_o := open_o.QueryRow(s)
   var artist_n int
   e = query_o.Scan(&artist_n)
   if e != nil {
      log.Fatal(e)
   }
   println(artist_n)
}
