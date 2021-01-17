package musicbrainz

import (
   "net/url"
   "sort"
   "winter"
)

const api = "https://musicbrainz.org/ws/2/release"

type MB struct {
   ID string
   Query url.Values
}

func New(mbid string) MB {
   query := url.Values{}
   query.Set("fmt", "json")
   query.Set("inc", "artist-credits recordings")
   return MB{mbid, query}
}

func (this MB) Group() (winter.Slice, error) {
   this.Query.Set("release-group", this.ID)
   url := api + "?" + this.Query.Encode()
   m, e := winter.JsonGetHttp(url)
   if e != nil {
      return nil, e
   }
   return m.A("releases"), nil
}

func (this MB) Release() (winter.Map, error) {
   url := api + "/" + this.ID + "?" + this.Query.Encode()
   return winter.JsonGetHttp(url)
}

func Status(m winter.Map) int {
   if m["status"] == nil {
      return 0
   }
   if m.S("status") != "Official" {
      return 0
   }
   return 1
}

func Date(m winter.Map, width int) string {
   left := ""
   if m["date"] != nil {
      left = m.S("date")
   }
   start := len(left)
   right := "9999-12-31"[start:]
   return (left + right)[:width]
}

func TrackLen(m winter.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
}

func Sort(a winter.Slice) {
   sort.Slice(a, func (first, second int) bool {
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
