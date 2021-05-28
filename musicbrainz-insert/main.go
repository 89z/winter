package main

import (
   "fmt"
   "os"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      return
   }
   db, err := sql.Open("sqlite3", os.Getenv("WINTER"))
   if err != nil {
      panic(err)
   }
   defer db.Close()
   tx, err := db.Begin()
   if err != nil {
      panic(err)
   }
   defer tx.Commit()
   r, err := release(os.Args[1])
   if err != nil {
      panic(err)
   }
   if err := insert(r, tx); err != nil {
      panic(err)
   }
}
