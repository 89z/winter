package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

type mbRelease struct {
   ReleaseCount int `json:"release-count"`
   Releases []struct {
      Date string
      Group struct {
         FirstRelease string `json:"first-release-date"`
         Id string
         SecondaryTypes []string `json:"secondary-types"`
         Title string
      } `json:"release-group"`
      Title string
   }
}

type values struct {
   offset int
   url.Values
}

func newValues(id string) values {
   value := values{}
   value.Set("fmt", "json")
   value.Set("inc", "release-groups")
   value.Set("limit", "100")
   value.Set("status", "official")
   value.Set("type", "album")
   value.Set("artist", id)
   return value
}

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release titles.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type winterRemote struct {
   color string
   date string
   release map[string]bool
   title string
}

func remoteAlbum(id string) ([]winterRemote, error) {
   var (
      remote = map[string]int{}
      remotes []winterRemote
      value = newValues(id)
   )
   for {
      get, e := http.Get(
         "http://musicbrainz.org/ws/2/release?" + value.Encode(),
      )
      if e != nil {
         return nil, e
      }
      var mb mbRelease
      e = json.NewDecoder(get.Body).Decode(&mb)
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
         index, ok := remote[release.Group.Id]
         if ok {
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
      value.offset += 100
      if value.offset >= mb.ReleaseCount {
         break
      }
      value.Set(
         "offset", fmt.Sprint(value.offset),
      )
   }
   return remotes, nil
}
