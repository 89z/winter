package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "winter/snow"
)

var offset_n float64

func RemoteAlbum(mb_s string) ([]Group, error) {
   q := url.Values{}
   q.Set("artist", mb_s)
   q.Set("fmt", "json")
   q.Set("inc", "release-groups")
   q.Set("limit", "100")
   q.Set("status", "official")
   q.Set("type", "album")
   remote_a, remote_m := []Group{}, map[string]int{}
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
         id_s := group_m.S("id")
         index_n, b := remote_m[id_s]
         if b {
            // add release to group
            remote_a[index_n].Release = append(
               remote_a[index_n].Release, release_m.S("title"),
            )
         } else {
            // add group
            remote_a = append(remote_a, Group{
               Date: group_m.S("first-release-date"),
               Release: []string{
                  release_m.S("title"),
               },
               Title: group_m.S("title"),
            })
            remote_m[id_s] = len(remote_a) - 1
         }
      }
      offset_n += 100
      if offset_n >= json_m.N("release-count") {
         break
      }
      q.Set("offset", fmt.Sprint(offset_n))
   }
   return remote_a, nil
}
