package musicbrainz

import (
   "encoding/json"
   "net/http"
   "net/url"
   "winter/assert"
)

type MB struct {
   API string
   ID string
   Query url.Values
}

func New(mbid_s string) MB {
   return MB{
      "https://musicbrainz.org/ws/2/release",
      mbid_s,
      url.Values{
         "fmt": []string{"json"},
         "inc": []string{"artist-credits recordings"},
      },
   }
}

func (this MB) Group() (assert.Slice, error) {
   this.Query["release-group"] = []string{this.ID}
   url_s := this.API + "?" + this.Query.Encode()
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
   url_s := this.API + "/" + this.ID + "?" + this.Query.Encode()
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
   return m, nil
}

func Date(m assert.Map) string {
   s := m.S("date")
   switch len(s) {
   case 0:
      return "9999-12-31"
   case 4:
      return s + "-12-31"
   case 6:
      return s + "-31"
   default:
      return s
   }
}

func IsOfficial(m assert.Map) bool {
   return m.S("status") == "Official"
}

func TrackLen(m assert.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
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
