package main

import (
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
      println(`Usage:
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
   album_m := map[string]string{"@date": rel_m.S("date")}
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
         album_m[track_m.S("title")] = note_s
      }
   }
   for key_s, val_s := range album_m {
      println(key_s, "|", val_s)
   }
}
