package main

import (
   "os"
   "path"
   "strings"
   "time"
   "winter"
   "winter/musicbrainz"
   "winter/youtube"
)

func main() {
   if len(os.Args) != 2 {
      println(`usage:
musicbrainz-views <URL>

examples:
https://musicbrainz.org/release-group/d03bb6b1-d7b4-38ea-974e-847cbb31dca4
https://musicbrainz.org/release/7a629d52-6a61-3ea1-a0a0-dd50bdef63b4`)
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
   out_a := []string{}
   artist_a := rel_m.A("artist-credit")
   for n := range artist_a {
      artist_s := artist_a.M(n).S("name")
      out_a = append(out_a, artist_s)
   }
   artist_s := strings.Join(out_a, " ")
   media_a := rel_m.A("media")
   for n := range media_a {
      track_a := media_a.M(n).A("tracks")
      for n := range track_a {
         title_s := track_a.M(n).S("title")
         id, e := youtubeResult(artist_s + " " + title_s)
         check(e)
         info_m, e := youtube.Info(id)
         check(e)
         view_n, e := youtube.Views(info_m)
         check(e)
         color_s, b := youtube.Color(view_n)
         println(color_s)
         if b {
            print("youtube.com/watch?v=", id, "\n")
            return
         }
         time.Sleep(500 * time.Millisecond)
      }
   }
}
