package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
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

Update artist pop:
   winter pop 999 0

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
      artist_s := os.Args[2]
      e = InsertArtist(open_o, artist_s)
   case "check":
      artist_s, check_s := os.Args[2], os.Args[3]
      e = UpdateCheck(open_o, artist_s, check_s)
   case "date":
      album_s, date_s := os.Args[2], os.Args[3]
      e = UpdateDate(open_o, album_s, date_s)
   case "note":
      song_s, note_s := os.Args[2], os.Args[3]
      e = UpdateNote(open_o, song_s, note_s)
   case "pop":
      artist_s, pop_s := os.Args[2], os.Args[3]
      e = UpdatePop(open_o, artist_s, pop_s)
   case "url":
      album_s, url_s := os.Args[2], os.Args[3]
      e = UpdateURL(open_o, album_s, url_s)
   default:
      e = SelectArtist(open_o, key_s)
   }
   if e != nil {
      log.Fatal(e)
   }
}
