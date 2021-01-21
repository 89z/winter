package musicbrainz

import (
   "github.com/89z/x"
   "github.com/89z/x/json"
   "net/url"
   "sort"
)

func date(m x.Map, width int) string {
   left := ""
   if m["date"] != nil {
      left = m.S("date")
   }
   start := len(left)
   right := "9999-12-31"[start:]
   return (left + right)[:width]
}

func Group(id string) (x.Slice, error) {
   q := url.Values{}
   q.Set("fmt", "json")
   q.Set("inc", "artist-credits recordings")
   q.Set("release-group", id)
   url := "https://musicbrainz.org/ws/2/release?" + q.Encode()
   m, e := json.LoadHttp(url)
   if e != nil {
      return nil, e
   }
   return m.A("releases"), nil
}

func Release(id string) (x.Map, error) {
   q := url.Values{}
   q.Set("fmt", "json")
   q.Set("inc", "artist-credits recordings")
   url := "https://musicbrainz.org/ws/2/release/" + id + "?" + q.Encode()
   return json.LoadHttp(url)
}

func Sort(a x.Slice) {
   sort.Slice(a, func (first, second int) bool {
      first_m, second_m := a.M(first), a.M(second)
      // 1. STATUS
      if status(first_m) > status(second_m) {
         return true
      }
      if status(first_m) < status(second_m) {
         return false
      }
      // 2. YEAR
      if date(first_m, 4) < date(second_m, 4) {
         return true
      }
      if date(first_m, 4) > date(second_m, 4) {
         return false
      }
      // 3. TRACKS
      if trackLen(first_m) < trackLen(second_m) {
         return true
      }
      if trackLen(first_m) > trackLen(second_m) {
         return false
      }
      // 4. DATE
      return date(first_m, 10) < date(second_m, 10)
   })
}

func status(m x.Map) int {
   if m["status"] == nil {
      return 0
   }
   if m.S("status") != "Official" {
      return 0
   }
   return 1
}

func trackLen(m x.Map) float64 {
   var track_n float64
   a := m.A("media")
   for n := range a {
      track_n += a.M(n).N("track-count")
   }
   return track_n
}
