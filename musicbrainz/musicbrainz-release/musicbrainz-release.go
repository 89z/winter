package main

import (
   "fmt"
   "log"
   "os"
   "path"
   "strings"
   "time"
   "winter/assert"
   "winter/musicbrainz"
)

const (
   max_n = 15 * time.Minute
   min_n = 179_500 * time.Millisecond
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`Usage:
musicbrainz-release <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   url_s := os.Args[1]
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
   open_o, e := sql.Open("sqlite3", "winter.db")
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   album_s := rel_m.S("title")
   date_s := rel_m.S("date")
   exec_o, e := open_o.Exec(
      "insert into album_t (album_s, date_s) values (?, ?)", album_s, date_s,
   )
   if e != nil {
      log.Fatal(e)
   }
   album_n, e := exec_o.LastInsertId()
   if e != nil {
      log.Fatal(e)
   }
   // ARTIST
   query_o := open_o.QueryRow(
      "select artist_n from artist_t where artist_s = ?",
      rel_m.A("artist-credit").M(0).S("name"),
   )
   var artist_n int
   e = query_o.Scan(&artist_n)
   if e != nil {
      log.Fatal(e)
   }
   // SONGS
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
   for song_s, note_s := range song_m {
      // SONG
      exec_o, e := open_o.Exec(
         "insert into song_t (song_s, note_s) values (?, ?)", song_s, note_s
      )
      if e != nil {
         log.Fatal(e)
      }
      song_n, e := exec_o.LastInsertId()
      if e != nil {
         log.Fatal(e)
      }
      // SONG ALBUM
      _, e := open_o.Exec(
         "insert into song_album_t values (?, ?)", song_n, album_n
      )
      if e != nil {
         log.Fatal(e)
      }
      // SONG ARTIST
      _, e := open_o.Exec(
         "insert into song_artist_t values (?, ?)", song_n, artist_n
      )
      if e != nil {
         log.Fatal(e)
      }
   }
}
