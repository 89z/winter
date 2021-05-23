package main

import (
   "database/sql"
   "os"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) == 1 {
      println(`Copy album:
   winter album 999 1000

Delete album:
   winter album 999

Select all artist:
   winter artist

Select one artist:
   winter 'Kate Bush'

Insert artist:
   winter artist 'Kate Bush'

Update artist date:
   winter check 999 2020-12-31

Update artist id:
   winter mb 999 3f5be744-e867-42fb-8913-5fd69e4099b5

Update album date:
   winter date 999 2020-12-31

Update album URL:
   winter url 999 youtube.com/watch?v=HQmmM_qwG4k

Update song note:
   winter note 999 good`)
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
   if err := transact(tx, os.Args[1], os.Args[2:]); err != nil {
      panic(err)
   }
}
