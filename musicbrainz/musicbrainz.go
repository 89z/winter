package musicbrainz

import (
   "github.com/89z/x"
   "github.com/89z/x/json"
   "net/url"
   "path"
   "sort"
   "strings"
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

func Release(in string) (x.Map, error) {
   id := path.Base(in)
   v := url.Values{}
   v.Set("fmt", "json")
   v.Set("inc", "artist-credits recordings")
   if strings.Contains(in, "/release/") {
      out := "https://musicbrainz.org/ws/2/release/" + id + "?" + v.Encode()
      return json.LoadHttp(out)
   }
   v.Set("release-group", id)
   out := "https://musicbrainz.org/ws/2/release?" + v.Encode()
   group, e := json.LoadHttp(out)
   if e != nil {
      return nil, e
   }
   albums := group.A("releases")
   Sort(albums)
   return albums.M(0), nil
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
