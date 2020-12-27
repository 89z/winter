package musicbrainz

import (
   "encoding/json"
   "net/http"
   "net/url"
   "winter/snow"
)

type MB struct {
   ID string
   Query url.Values
}

const API = "https://musicbrainz.org/ws/2/release"

func New(mbid_s string) MB {
   m := url.Values{}
   m.Set("fmt", "json")
   m.Set("inc", "artist-credits recordings")
   return MB{mbid_s, m}
}

func (this MB) Group() (snow.Slice, error) {
   this.Query.Set("release-group", this.ID)
   url_s := API + "?" + this.Query.Encode()
   println(url_s)
   o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   json_m := snow.Map{}
   e = json.NewDecoder(o.Body).Decode(&json_m)
   if e != nil {
      return nil, e
   }
   release_a := json_m.A("releases")
   official_a := snow.Slice{}
   for n := range release_a {
      release_m := release_a.M(n)
      if release_m["status"] == nil {
         continue
      }
      if release_m.S("status") != "Official" {
         continue
      }
      official_a = append(official_a, release_m)
   }
   return official_a, nil
}

func (this MB) Release() (snow.Map, error) {
   url_s := API + "/" + this.ID + "?" + this.Query.Encode()
   println(url_s)
   o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   m := snow.Map{}
   return m, json.NewDecoder(o.Body).Decode(&m)
}

func TrackLen(m snow.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
}

func Date(m snow.Map, width int) string {
   left := ""
   if m["date"] != nil {
      left = m.S("date")
   }
   start := len(left)
   right := "9999-12-31"[start:]
   return (left + right)[:width]
}

func Reduce(acc_n int, cur_m snow.Map, cur_n int, src_a snow.Slice) int {
   if cur_n == 0 {
      return 0
   }
   acc_m := src_a.M(acc_n)
   // 1. YEAR
   if Date(cur_m, 4) > Date(acc_m, 4) {
      return acc_n
   }
   if Date(cur_m, 4) < Date(acc_m, 4) {
      return cur_n
   }
   // 2. TRACKS
   if TrackLen(cur_m) > TrackLen(acc_m) {
      return acc_n
   }
   if TrackLen(cur_m) < TrackLen(acc_m) {
      return cur_n
   }
   // 3. DATE
   if Date(cur_m, 10) >= Date(acc_m, 10) {
      return acc_n
   }
   return cur_n
}
