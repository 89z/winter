package main

import (
   "fmt"
   "github.com/89z/sienna/json"
   "net/url"
)

var offset float64

func remoteAlbum(mb_s string) ([]group, error) {
   q := url.Values{}
   q.Set("artist", mb_s)
   q.Set("fmt", "json")
   q.Set("inc", "release-groups")
   q.Set("limit", "100")
   q.Set("status", "official")
   q.Set("type", "album")
   remote_a, remote_m := []group{}, map[string]int{}
   for {
      url := "https://musicbrainz.org/ws/2/release?" + q.Encode()
      mb, e := json.LoadHttp(url)
      if e != nil {
         return nil, e
      }
      release_a := mb.A("releases")
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
         id := group_m.S("id")
         index_n, b := remote_m[id]
         release_s := release_m.S("title")
         if b {
            // add release to group
            remote_a[index_n].release[release_s] = true
         } else {
            // add group
            remote_a = append(remote_a, group{
               date: group_m.S("first-release-date"),
               release: map[string]bool{release_s: true},
               title: group_m.S("title"),
            })
            remote_m[id] = len(remote_a) - 1
         }
      }
      offset += 100
      if offset >= mb.N("release-count") {
         break
      }
      q.Set("offset", fmt.Sprint(offset))
   }
   return remote_a, nil
}

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release titles.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type group struct {
   color string
   date string
   release map[string]bool
   title string
}
