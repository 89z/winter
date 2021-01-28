package main

import (
   "fmt"
   "github.com/89z/x/json"
   "net/url"
)

var offset float64

func remoteAlbum(mb_s string) ([]winterGroup, error) {
   value := make(url.Values)
   value.Set("artist", mb_s)
   value.Set("fmt", "json")
   value.Set("inc", "release-groups")
   value.Set("limit", "100")
   value.Set("status", "official")
   value.Set("type", "album")
   remotes, remote := []winterGroup{}, map[string]int{}
   for {
      get, e := http.Get(
         "https://musicbrainz.org/ws/2/release?" + value.Encode(),
      )
      if e != nil {
         return nil, e
      }
      mb := new(mbRelease)
      e = json.NewDecoder(get.Body).Decode(mb)
      if e != nil {
         return nil, e
      }
      releases := mb.A("releases")
      for n := range releases {
         release := releases.M(n)
         group := release.M("release-group")
         if release["date"] == nil {
            continue
         }
         if release.S("date") == "" {
            continue
         }
         if len(group.A("secondary-types")) > 0 {
            continue
         }
         id := group.S("id")
         release_s := release.S("title")
         index_n, b := remote[id]
         if b {
            // add release to group
            remotes[index_n].release[release_s] = true
         } else {
            // add group
            remotes = append(remotes, winterGroup{
               date: group.S("first-release-date"),
               release: map[string]bool{release_s: true},
               title: group.S("title"),
            })
            remote[id] = len(remotes) - 1
         }
      }
      offset += 100
      if offset >= mb.N("release-count") {
         break
      }
      value.Set("offset", fmt.Sprint(offset))
   }
   return remotes, nil
}
