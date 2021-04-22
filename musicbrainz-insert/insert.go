package main

import (
   "fmt"
   "github.com/89z/x/musicbrainz"
   "log"
   "os"
   "time"
   "winter"
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

type titleNote struct { title, note string }

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   tx, e := winter.NewTx(os.Getenv("WINTER"))
   if e != nil {
      log.Fatal(e)
   }
   album, e := musicbrainz.NewRelease(os.Args[1])
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   albumId, e := tx.Insert(
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album.Title,
      album.Date,
   )
   if e != nil {
      log.Fatal(e)
   }
   // CREATE SONG ARRAY
   var songs []titleNote
   for _, media := range album.Media {
      for _, track := range media.Tracks {
         songs = append(songs, titleNote{
            track.Title, note(track.Length),
         })
      }
   }
   // CREATE ARTIST ARRAY
   var artists []int
   for _, each := range album.ArtistCredit {
      var artist int
      e = tx.QueryRow(
         "select artist_n from artist_t where mb_s = ?", each.Artist.Id,
      ).Scan(&artist)
      if e != nil {
         log.Fatalln(each.Name, e)
      }
      artists = append(artists, artist)
   }
   // ITERATE SONG ARRAY
   for _, each := range songs {
      song, e := tx.Insert(
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         each.title,
         each.note,
         albumId,
      )
      if e != nil {
         log.Fatal(e)
      }
      // ITERATE ARTIST ARRAY
      for _, artist := range artists {
         _, e = tx.Insert("song_artist_t values (?, ?)", song, artist)
         if e != nil {
            log.Fatal(e)
         }
      }
   }
   e = tx.Commit()
   if e != nil {
      log.Fatal(e)
   }
}
