package main

import (
   "net/http"
   "net/url"
   "os"
   "path"
   "strings"
   "winter/musicbrainz"
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
   rel_m := assert.Map{}
   if strings.Contains(url_s, "release-group") {
      rel_a := mb_o.Group()
      rel_n := 0
      for idx_n := range rel_a {
         cur_m := rel_a.M(idx_n)
         rel_n = musicbrainz.Reduce(rel_n, cur_m, idx_n, rel_a)
      }
      rel_m = rel_a.M(rel_n)
      print("musicbrainz.org/release/", rel_m.S("id"), "\n")
   } else {
      rel_m = mb_o.Release()
   }
   out_a := assert.Slice{}
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
         id_s, e := YoutubeResult(artist_s + " " + title_s)
         if e != nil {
            log.Fatal(e)
         }
         $info_o = youtube_info($id_s);
         echo youtube_views($info_o), "\n\n";
         usleep(500_000);
      }
   }
}
