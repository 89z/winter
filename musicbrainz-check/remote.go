package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

var offset int

func remoteAlbum(mb_s string) ([]winterRemote, error) {
   value := make(url.Values)
   value.Set("artist", mb_s)
   value.Set("fmt", "json")
   value.Set("inc", "release-groups")
   value.Set("limit", "100")
   value.Set("status", "official")
   value.Set("type", "album")
   remotes, remote := []winterRemote{}, map[string]int{}
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
      for _, release := range mb.Releases {
         if release.Date == "" {
            continue
         }
         if len(release.Group.SecondaryTypes) > 0 {
            continue
         }
         index, b := remote[release.Group.Id]
         if b {
            // add release to group
            remotes[index].release[release.Title] = true
         } else {
            // add group
            remotes = append(remotes, winterRemote{
               date: release.Group.FirstRelease,
               release: map[string]bool{release.Title: true},
               title: release.Group.Title,
            })
            remote[release.Group.Id] = len(remotes) - 1
         }
      }
      offset += 100
      if offset >= mb.ReleaseCount {
         break
      }
      value.Set("offset", fmt.Sprint(offset))
   }
   return remotes, nil
}
