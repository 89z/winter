package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
   "winter/snow"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) == 1 {
      fmt.Println(`winter <target> <arguments>

Delete album:
   winter album 999

Select artist:
   winter 'Kate Bush'

Insert artist:
   winter artist 'Kate Bush'

Update artist date:
   winter check 999 2019-12-31

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
      album_s := os.Args[2]
      e = DeleteAlbum(open_o, album_s)
   case "artist":
      _, e = snow.Insert(
         open_o,
         "artist_t (artist_s, check_s) values (?, '')",
         os.Args[2],
      )
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
      e = SelectArtist(open_o, key_s)
   }
   if e != nil {
      log.Fatal(e)
   }
}
