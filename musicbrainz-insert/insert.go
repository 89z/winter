package main

import (
   "database/sql"
   "fmt"
   "github.com/89z/winter"
   "github.com/89z/x"
   "github.com/89z/x/musicbrainz"
   "os"
   "time"
   _ "github.com/mattn/go-sqlite3"
)

type titleNote struct {
   note, title string
}

func begin(file string) (*sql.Tx, error) {
   db, e := sql.Open("sqlite3", file)
   if e != nil {
      return e
   }
   return db.Begin()
}

func note(length int) string {
   if length == 0 {
      return "?:??"
   }
   dur := time.Duration(length) * time.Millisecond
   if dur < 179_500 * time.Millisecond {
      return "short"
   }
   if dur > 15 * time.Minute {
      return "long"
   }
   return ""
}

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   tx, e := begin(
      os.Getenv("WINTER"),
   )
   x.Check(e)
   album, e := musicbrainz.NewRelease(os.Args[1])
   x.Check(e)
   // CREATE ARTIST ARRAY
   var artists []int
   for _, each := range album.ArtistCredit {
      var artist int
      e = tx.QueryRow(
         "select artist_n from artist_t where mb_s = ?", each.Artist.Id,
      ).Scan(&artist)
      x.Check(e, each.Name)
      artists = append(artists, artist)
   }
   // ALBUM
   albumId, e := winter.Insert(
      tx,
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album.Title,
      album.Date,
   )
   x.Check(e)
   // CREATE SONG ARRAY
   var songs []titleNote
   for _, media := range album.Media {
      for _, track := range media.Tracks {
         songs = append(songs, titleNote{
            track.Title, note(track.Length),
         })
      }
   }
   // ITERATE SONG ARRAY
   for _, each := range songs {
      song, e := winter.Insert(
         tx,
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         each.title,
         each.note,
         albumId,
      )
      x.Check(e)
      // ITERATE ARTIST ARRAY
      for _, artist := range artists {
         _, e = winter.Insert(tx, "song_artist_t values (?, ?)", song, artist)
         x.Check(e)
      }
   }
   e = tx.Commit()
   x.Check(e)
}
