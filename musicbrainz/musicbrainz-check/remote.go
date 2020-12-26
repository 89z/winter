package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "winter/snow"
)

type Remote struct {
   Date string
   Group string
   Title string
}

var (
   offset_n float64
   remote_a []Remote
)

func RemoteAlbum(mb_s string) ([]Remote, error) {
   q := url.Values{}
   q.Set("artist", mb_s)
   q.Set("fmt", "json")
   q.Set("inc", "release-groups")
   q.Set("limit", "100")
   q.Set("status", "official")
   q.Set("type", "album")
   for {
      url_s := "https://musicbrainz.org/ws/2/release?" + q.Encode()
      fmt.Println(url_s)
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
      for n := range release_a {
         release_m := release_a.M(n)
         group_m := release_m.M("release-group")
         if release_m["date"] == nil {
            continue
         }
         date_s := release_m.S("date")
         if date_s == "" {
            continue
         }
         second_a := group_m.A("secondary-types")
         if len(second_a) > 0 {
            continue
         }
         remote_a = append(remote_a, Remote{
            date_s, group_m.S("id"), release_m.S("title"),
         })
      }
      offset_n += 100
      if offset_n >= json_m.N("release-count") {
         break
      }
      q.Set("offset", fmt.Sprint(offset_n))
   }
   return remote_a, nil
}
