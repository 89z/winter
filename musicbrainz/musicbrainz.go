package musicbrainz

import (
   "encoding/json"
   "net/http"
   "net/url"
   "winter/snow"
   "sort"
)

const API = "https://musicbrainz.org/ws/2/release"

type MB struct {
   ID string
   Query url.Values
}

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
   return json_m.A("releases"), nil
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

func Status(m snow.Map) int {
   if m["status"] == nil {
      return 0
   }
   if m.S("status") != "Official" {
      return 0
   }
   return 1
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

func TrackLen(m snow.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
}

func Sort(a snow.Slice) {
   sort.Slice(a, func(first, second int) bool {
      first_m, second_m := a.M(first), a.M(second)
      // 1. STATUS
      if Status(first_m) > Status(second_m) {
         return true
      }
      if Status(first_m) < Status(second_m) {
         return false
      }
      // 2. YEAR
      if Date(first_m, 4) < Date(second_m, 4) {
         return true
      }
      if Date(first_m, 4) > Date(second_m, 4) {
         return false
      }
      // 3. TRACKS
      if TrackLen(first_m) < TrackLen(second_m) {
         return true
      }
      if TrackLen(first_m) > TrackLen(second_m) {
         return false
      }
      // 4. DATE
      return Date(first_m, 10) < Date(second_m, 10)
   })
}
