package main

import (
   "database/sql"
   "flag"
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
   var confirm_b bool
   flag.BoolVar(&confirm_b, "c", false, "confirm")
   flag.Parse()
   if flag.NArg() != 1 {
      fmt.Println(`musicbrainz-insert [flags] <URL>

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
   rel_m := snow.Map{}
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
   songs_a := [][]string{}
   media_a := rel_m.A("media")
   for n := range media_a {
      track_a := media_a.M(n).A("tracks")
      for n := range track_a {
         track_m := track_a.M(n)
         song_s := track_m.S("title")
         note_s := Note(track_m.N("length"))
         songs_a = append(songs_a, []string{song_s, note_s})
      }
   }
   if ! confirm_b {
      fmt.Printf("artist: %q\n", artist_s)
      fmt.Printf("album: %q\n", album_s)
      fmt.Printf("date: %q\n", date_s)
      for _, song_a := range songs_a {
         fmt.Printf("song: %q, note: %q\n", song_a[0], song_a[1])
      }
      return
   }
   db_s := os.Getenv("WINTER")
   open_o, e := sql.Open("sqlite3", db_s)
   if e != nil {
      log.Fatal(e)
   }
   // ALBUM
   album_n, e := snow.Insert(
      open_o,
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
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
   for _, song_a := range songs_a {
      // SONG
      song_n, e := snow.Insert(
         open_o, "song_t (song_s, note_s) values (?, ?)", song_a[0], song_a[1],
      )
      if e != nil {
         log.Fatal(e)
      }
      // SONG ALBUM
      _, e = snow.Insert(
         open_o, "song_album_t values (?, ?)", song_n, album_n,
      )
      if e != nil {
         log.Fatal(e)
      }
      // SONG ARTIST
      _, e = snow.Insert(
         open_o, "song_artist_t values (?, ?)", song_n, artist_n,
      )
      if e != nil {
         log.Fatal(e)
      }
   }
}
