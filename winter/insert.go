package main

import (
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

func InsertArtist(open_o *sql.DB, artist_s string) error {
   _, e := open_o.Exec(
      "insert into artist_t (artist_s, check_s, pop_n) values (?, '', 1)",
      artist_s,
   )
   return e
}
