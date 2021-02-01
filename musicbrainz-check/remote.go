package main

import (
   "encoding/json"
   "fmt"
   "github.com/89z/x"
   "net/http"
)

var (
   offset int
   remote = map[string]int{}
   remotes []winterRemote
)

func remoteAlbum(id string) ([]winterRemote, error) {
   url := x.NewURL()
   url.Host = "musicbrainz.org"
   url.Path = "ws/2/release"
   url.QuerySet("fmt", "json")
   url.QuerySet("inc", "release-groups")
   url.QuerySet("limit", "100")
   url.QuerySet("status", "official")
   url.QuerySet("type", "album")
   url.QuerySet("artist", id)
   for {
      get, e := http.Get(
         url.String(),
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
      url.QuerySet(
         "offset", fmt.Sprint(offset),
      )
   }
   return remotes, nil
}
