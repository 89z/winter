package main

import (
   "database/sql"
   "fmt"
   "github.com/89z/winter"
   "github.com/89z/x/musicbrainz"
   "log"
   "os"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   url := os.Args[1]
   album, e := musicbrainz.Release(url)
   if e != nil {
      log.Fatal(e)
   }
   album_s := album.S("title")
   date_s := album.S("date")
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   if e != nil {
      log.Fatal(e)
   }
   tx, e := db.Begin()
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   album_n, e := winter.Insert(
      tx,
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album_s,
      date_s,
   )
   if e != nil {
      log.Fatal(e)
   }
   var (
      artist_n int
      artists []int
      songs []song
   )
   // CREATE ARTIST ARRAY
   credit := album.A("artist-credit")
   for n := range credit {
      // Chicago, Chicago Transit Authority
      name := credit.M(n).M("artist").S("name")
      query := tx.QueryRow(
         "select artist_n from artist_t where artist_s = ?", name,
      )
      e = query.Scan(&artist_n)
      if e != nil {
         log.Fatalln(name, e)
      }
      artists = append(artists, artist_n)
   }
   // CREATE SONG ARRAY
   media := album.A("media")
   for n := range media {
      track_a := media.M(n).A("tracks")
      for n := range track_a {
         track_m := track_a.M(n)
         song_s := track_m.S("title")
         note_s := note(track_m)
         songs = append(songs, song{song_s, note_s})
      }
   }
   // ITERATE SONG ARRAY
   for _, song_o := range songs {
      song_n, e := winter.Insert(
         tx,
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         song_o.title,
         song_o.note,
         album_n,
      )
      if e != nil {
         log.Fatal(e)
      }
      // ITERATE ARTIST ARRAY
      for _, artist_n := range artists {
         _, e = winter.Insert(
            tx, "song_artist_t values (?, ?)", song_n, artist_n,
         )
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
