package musicbrainz

import (
   "encoding/json"
   "net/http"
   "net/url"
   "winter/assert"
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

func (this MB) Group() (assert.Slice, error) {
   this.Query.Set("release-group", this.ID)
   url_s := API + "?" + this.Query.Encode()
   println(url_s)
   o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   m := assert.Map{}
   e = json.NewDecoder(o.Body).Decode(&m)
   if e != nil {
      return nil, e
   }
   return m.A("releases"), nil
}

func (this MB) Release() (assert.Map, error) {
   url_s := API + "/" + this.ID + "?" + this.Query.Encode()
   println(url_s)
   o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   m := assert.Map{}
   return m, json.NewDecoder(o.Body).Decode(&m)
}

func IsOfficial(m assert.Map) bool {
   return m["status"] != nil && m.S("status") == "Official"
}

func TrackLen(m assert.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
}

func Date(m assert.Map) string {
   s := ""
   if m["date"] != nil {
      s = m.S("date")
   }
   n := len(s)
   return s + "9999-12-31"[n:]
}

func Reduce(old_n int, new_m assert.Map, new_n int, old_a assert.Slice) int {
   if new_n == 0 {
      return 0
   }
   old_m := old_a.M(old_n)
   if ! IsOfficial(new_m) {
      return old_n
   }
   if Date(new_m) > Date(old_m) {
      return old_n
   }
   if Date(new_m) < Date(old_m) {
      return new_n
   }
   if TrackLen(new_m) >= TrackLen(old_m) {
      return old_n
   }
   return new_n
}
