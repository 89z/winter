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
   // ARTIST
   artist_s := rel_m.A("artist-credit").M(0).S("name")
   fmt.Println("artist_t:")
   fmt.Printf("\tartist_s: %q\n", artist_s)
   fmt.Println("\tartist_n: ?")
   // ALBUM
   album_s := rel_m.S("title")
   date_s := rel_m.S("date")
   fmt.Println("album_t:")
   fmt.Println("\talbum_n: ?")
   fmt.Printf("\talbum_s: %q\n", album_s)
   fmt.Printf("\tdate_s: %q\n", date_s)
   fmt.Printf("\turl_s: %q\n", "")
   media_a := rel_m.A("media")
   for n := range media_a {
      track_a := media_a.M(n).A("tracks")
      for n := range track_a {
         track_m := track_a.M(n)
         len_n := time.Duration(track_m.N("length")) * time.Millisecond
         note_s := ""
         if len_n < min_n {
            note_s = "short"
         }
         if len_n > max_n {
            note_s = "long"
         }
         fmt.Println("--------------------------------------------------------")
         // SONG
         title_s := track_m.S("title")
         fmt.Println("song_t:")
         fmt.Printf("\tsong_n: ?, song_s: %q, note_s: %q\n", title_s, note_s)
         // SONG ALBUM
         fmt.Println("song_album_t:")
         fmt.Println("\tsong_n: ?, album_n: ?")
         // SONG ARTIST
         fmt.Println("song_artist_t:")
         fmt.Println("\tsong_n: ?, artist_n: ?")
      }
   }
}
