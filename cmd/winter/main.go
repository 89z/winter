package main

import (
   "database/sql"
   "log"
   "os"
   "winter"
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
      os.Exit(1)
   }
   key_s := os.Args[1]
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   if e != nil {
      log.Fatal(e)
   }
   tx, e := db.Begin()
   if e != nil {
      log.Fatal(e)
   }
   switch key_s {
   case "album":
      source := os.Args[2]
      if len(os.Args) == 4 {
         dest := os.Args[3]
         e = copyAlbum(tx, source, dest)
      } else {
         e = deleteAlbum(tx, source)
      }
   case "artist":
      if len(os.Args) == 2 {
         e = selectAll(tx)
      } else {
         _, e = winter.Insert(
            tx,
            "artist_t (artist_s, check_s, mb_s) values (?, '', '')",
            os.Args[2],
         )
      }
   case "check":
      e = winter.Update(
         tx,
         "artist_t set check_s = ? where artist_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "date":
      e = winter.Update(
         tx,
         "album_t set date_s = ? where album_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "mb":
      e = winter.Update(
         tx,
         "artist_t set mb_s = ? where artist_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "note":
      e = winter.Update(
         tx,
         "song_t set note_s = ? where song_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "url":
      e = winter.Update(
         tx,
         "album_t set url_s = ? where album_n = ?",
         os.Args[3],
         os.Args[2],
      )
   default:
      e = selectOne(tx, key_s)
   }
   if e != nil {
      log.Fatal(e)
   }
   e = tx.Commit()
   if e != nil {
      log.Fatal(e)
   }
}
