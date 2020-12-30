package main

import (
   "database/sql"
   "log"
   "os"
   "winter/snow"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) == 1 {
      println(`winter <target> <arguments>

Copy album:
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
   winter check 999 2019-12-31

Update artist id:
   winter mb 999 3f5be744-e867-42fb-8913-5fd69e4099b5

Update album date:
   winter date 999 2019-12-31

Update album URL:
   winter url 999 youtube.com/watch?v=HQmmM_qwG4k

Update song note:
   winter note 999 good`)
      os.Exit(1)
   }
   db_s := os.Getenv("WINTER")
   open_o, e := sql.Open("sqlite3", db_s)
   if e != nil {
      log.Fatal(e)
   }
   key_s := os.Args[1]
   switch key_s {
   case "album":
      source := os.Args[2]
      if len(os.Args) == 4 {
         dest := os.Args[3]
         e = CopyAlbum(open_o, source, dest)
      } else {
         e = DeleteAlbum(open_o, source)
      }
   case "artist":
      if len(os.Args) == 2 {
         e = SelectAll(open_o)
      } else {
         _, e = snow.Insert(
            open_o,
            "artist_t (artist_s, check_s, mb_s) values (?, '', '')",
            os.Args[2],
         )
      }
   case "check":
      e = snow.Update(
         open_o,
         "artist_t set check_s = ? where artist_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "date":
      e = snow.Update(
         open_o,
         "album_t set date_s = ? where album_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "mb":
      e = snow.Update(
         open_o,
         "artist_t set mb_s = ? where artist_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "note":
      e = snow.Update(
         open_o,
         "song_t set note_s = ? where song_n = ?",
         os.Args[3],
         os.Args[2],
      )
   case "url":
      e = snow.Update(
         open_o,
         "album_t set url_s = ? where album_n = ?",
         os.Args[3],
         os.Args[2],
      )
   default:
      e = SelectOne(open_o, key_s)
   }
   if e != nil {
      log.Fatal(e)
   }
}
