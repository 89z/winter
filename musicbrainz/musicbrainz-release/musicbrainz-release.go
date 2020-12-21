package main

import (
   "database/sql"
   "flag"
   "fmt"
   "log"
   "os"
   "path"
   "strings"
   "time"
   "winter/assert"
   "winter/musicbrainz"
   _ "github.com/mattn/go-sqlite3"
)

const (
   max_n = 15 * time.Minute
   min_n = 179_500 * time.Millisecond
)

func main() {
   var confirm_b bool
   flag.BoolVar(&confirm_b, "c", false, "confirm")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println(`musicbrainz-release [flags] <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872

flags:`)
      flag.PrintDefaults()
      os.Exit(1)
   }
   url_s := flag.Arg(0)
   mbid_s := path.Base(url_s)
   mb_o := musicbrainz.New(mbid_s)
   rel_m := assert.Map{}
   if strings.Contains(url_s, "release-group") {
      rel_a, e := mb_o.Group()
      if e != nil {
         log.Fatal(e)
      }
      rel_n := 0
      for idx_n := range rel_a {
         cur_m := rel_a.M(idx_n)
         rel_n = musicbrainz.Reduce(rel_n, cur_m, idx_n, rel_a)
      }
      rel_m = rel_a.M(rel_n)
   } else {
      var e error
      rel_m, e = mb_o.Release()
      if e != nil {
         log.Fatal(e)
      }
   }
   // Chicago, Chicago Transit Authority
   artist_s := rel_m.A("artist-credit").M(0).M("artist").S("name")
   album_s := rel_m.S("title")
   date_s := rel_m.S("date")
   song_m := map[string]string{}
   media_a := rel_m.A("media")
   for n := range media_a {
      track_a := media_a.M(n).A("tracks")
      for n := range track_a {
         track_m := track_a.M(n)
         song_s := track_m.S("title")
         len_n := time.Duration(track_m.N("length")) * time.Millisecond
         if len_n < min_n {
            song_m[song_s] = "short"
         } else if len_n > max_n {
            song_m[song_s] = "long"
         } else {
            song_m[song_s] = ""
         }
      }
   }
   if ! confirm_b {
      fmt.Println("artist_s:", artist_s)
      fmt.Println("album_s:", album_s)
      fmt.Println("date_s:", date_s)
      for song_s, note_s := range song_m {
         fmt.Print("song_s: ", song_s, ", note_s: ", note_s, "\n")
      }
      return
   }
   db_s := os.Getenv("WINTER")
   open_o, e := sql.Open("sqlite3", db_s)
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   album_n, e := Exec(
      open_o,
      "insert into album_t (album_s, date_s, url_s) values (?, ?, '')",
      album_s,
      date_s,
   )
   if e != nil {
      log.Fatal(e)
   }
   // ARTIST
   query_o := open_o.QueryRow(
      "select artist_n from artist_t where artist_s = ?", artist_s,
   )
   var artist_n int
   e = query_o.Scan(&artist_n)
   if e != nil {
      log.Fatal(e)
   }
   // SONGS
   for song_s, note_s := range song_m {
      // SONG
      song_n, e := Exec(
         open_o,
         "insert into song_t (song_s, note_s) values (?, ?)",
         song_s,
         note_s,
      )
      if e != nil {
         log.Fatal(e)
      }
      // SONG ALBUM
      _, e = Exec(
         open_o, "insert into song_album_t values (?, ?)", song_n, album_n,
      )
      if e != nil {
         log.Fatal(e)
      }
      // SONG ARTIST
      _, e = Exec(
         open_o, "insert into song_artist_t values (?, ?)", song_n, artist_n,
      )
      if e != nil {
         log.Fatal(e)
      }
   }
}
