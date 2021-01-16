package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
   "path"
   "strings"
   "winter"
   "winter/musicbrainz"
   _ "github.com/mattn/go-sqlite3"
)

func check(e error) {
   if e != nil {
      log.Fatal(e)
   }
}

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
   rel_m := winter.Map{}
   if strings.Contains(url_s, "release-group") {
      rel_a, e := mb_o.Group()
      check(e)
      musicbrainz.Sort(rel_a)
      rel_m = rel_a.M(0)
   } else {
      var e error
      rel_m, e = mb_o.Release()
      check(e)
   }
   album_s := rel_m.S("title")
   date_s := rel_m.S("date")
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   check(e)
   tx, e := db.Begin()
   check(e)
   // ALBUM
   album_n, e := winter.Insert(
      tx,
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album_s,
      date_s,
   )
   check(e)
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
      query_o := tx.QueryRow(
         "select artist_n from artist_t where artist_s = ?", name_s,
      )
      e = query_o.Scan(&artist_n)
      check(e)
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
      song_n, e := winter.Insert(
         tx,
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         song_o.Title,
         song_o.Note,
         album_n,
      )
      check(e)
      // ITERATE ARTIST ARRAY
      for _, artist_n := range artist_a {
         _, e = winter.Insert(
            tx, "song_artist_t values (?, ?)", song_n, artist_n,
         )
         check(e)
      }
   }
   e = tx.Commit()
   check(e)
}
