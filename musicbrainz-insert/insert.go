package main

import (
   "fmt"
   "github.com/89z/musicbrainz"
   "os"
   "time"
   "database/sql"
   _ "github.com/mattn/go-sqlite3"
)

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

type titleNote struct {
   title string
   note string
}

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
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
   album, err := musicbrainz.NewRelease(os.Args[1])
   if err != nil {
      panic(err)
   }
   err = insert(album, tx)
   if err != nil {
      panic(err)
   }
}

func insert(album musicbrainz.Release, tx *sql.Tx) error {
   // ALBUM
   albumId, err := tx.Exec(`
   INSERT INTO album_t (album_s, date_s, url_s) VALUES (?, ?, '')
   `, album.Title, album.Date)
   if err != nil { return err }
   // CREATE ARTIST ARRAY
   var artists []int
   for _, credit := range album.ArtistCredit {
      var artist int
      err := tx.QueryRow(`
      SELECT artist_n FROM artist_t WHERE mb_s = ?
      `, credit.Artist.ID).Scan(&artist)
      if err != nil {
         return fmt.Errorf("%v %v", credit.Name, err)
      }
      artists = append(artists, artist)
   }
   // CREATE SONG ARRAY
   var tns []titleNote
   for _, media := range album.Media {
      for _, track := range media.Tracks {
         tns = append(tns, titleNote{
            track.Title, note(track.Length),
         })
      }
   }
   // ITERATE SONG ARRAY
   for _, tn := range tns {
      song, err := tx.Exec(`
      INSERT INTO song_t (song_s, note_s, album_n) VALUES (?, ?, ?)
      `, tn.title, tn.note, albumId)
      if err != nil { return err }
      // ITERATE ARTIST ARRAY
      for _, artist := range artists {
         _, err := tx.Exec(`
         INSERT INTO song_artist_t VALUES (?, ?)
         `, song, artist)
         if err != nil { return err }
      }
   }
   return nil
}
