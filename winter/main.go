package main

import (
   "os"
   "winter"
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
   key := os.Args[1]
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
   switch key {
   case "album":
      source := os.Args[2]
      if len(os.Args) == 4 {
         err = copyAlbum(tx, source, os.Args[3])
      } else {
         err = deleteAlbum(tx, source)
      }
   case "artist":
      if len(os.Args) == 2 {
         err = selectAll(tx)
      } else {
         _, err = tx.Exec(`
         INSERT INTO artist_t (artist_s, check_s, mb_s) VALUES (?, '', '')
         `, os.Args[2])
      }
   case "check":
      err = tx.Exec(`
      UPDATE artist_t SET check_s = ? WHERE artist_n = ?
      `, os.Args[3], os.Args[2])
   case "date":
      err = tx.Exec(`
      UPDATE album_t SET date_s = ? WHERE album_n = ?
      `, os.Args[3], os.Args[2])
   case "mb":
      err = tx.Exec(`
      UPDATE artist_t SET mb_s = ? WHERE artist_n = ?
      `, os.Args[3], os.Args[2])
   case "note":
      err = tx.Exec(`
      UPDATE song_t SET note_s = ? WHERE song_n = ?
      `, os.Args[3], os.Args[3])
   case "url":
      err = tx.Exec(`
      UPDATE album_t SET url_s = ? WHERE album_n = ?
      `, os.Args[3], os.Args[2])
   default:
      err = selectOne(tx, key)
   }
   if err != nil {
      panic(err)
   }
}
