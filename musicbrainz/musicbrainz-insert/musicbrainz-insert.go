package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
   "path"
   "strings"
   "winter/musicbrainz"
   "winter/snow"
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
   url_s := os.Args[1]
   mbid_s := path.Base(url_s)
   mb_o := musicbrainz.New(mbid_s)
   rel_m := snow.Map{}
   if strings.Contains(url_s, "release-group") {
      rel_a, e := mb_o.Group()
      if e != nil {
         log.Fatal(e)
      }
      musicbrainz.Sort(rel_a)
      rel_m = rel_a.M(0)
   } else {
      var e error
      rel_m, e = mb_o.Release()
      if e != nil {
         log.Fatal(e)
      }
   }
   album_s := rel_m.S("title")
   date_s := rel_m.S("date")
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   album_n, e := snow.Insert(
      db,
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album_s,
      date_s,
   )
   if e != nil {
      log.Fatal(e)
   }
   var (
      artist_a []int
      artist_n int
      song_a []Song
   )
   // CREATE ARTIST ARRAY
   credit_a := rel_m.A("artist-credit")
   for n := range credit_a {
      // Chicago, Chicago Transit Authority
      name_s := credit_a.M(n).M("artist").S("name")
      query_o := db.QueryRow(
         "select artist_n from artist_t where artist_s = ?", name_s,
      )
      e = query_o.Scan(&artist_n)
      if e != nil {
         log.Fatal(e)
      }
      artist_a = append(artist_a, artist_n)
   }
   // CREATE SONG ARRAY
   media_a := rel_m.A("media")
   for n := range media_a {
      track_a := media_a.M(n).A("tracks")
      for n := range track_a {
         track_m := track_a.M(n)
         song_s := track_m.S("title")
         note_s := Note(track_m)
         song_a = append(song_a, Song{song_s, note_s})
      }
   }
   // ITERATE SONG ARRAY
   for _, song_o := range song_a {
      song_n, e := snow.Insert(
         db,
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         song_o.Title,
         song_o.Note,
         album_n,
      )
      if e != nil {
         log.Fatal(e)
      }
      // ITERATE ARTIST ARRAY
      for _, artist_n := range artist_a {
         _, e = snow.Insert(
            db, "song_artist_t values (?, ?)", song_n, artist_n,
         )
         if e != nil {
            log.Fatal(e)
         }
      }
   }
}
