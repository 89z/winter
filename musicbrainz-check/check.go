package main

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "winter"
)

type winterLocal struct {
   color string
   date string
}

func color(url string, unrated, good int) string {
   const (
      block = "\u2587\u2587\u2587\u2587\u2587"
      greenFive = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[m"
      greenTen = "\x1b[92m" + block + block + "\x1b[m"
      redFive = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[m"
      redTen = "\x1b[91m" + block + block + "\x1b[m"
   )
   if strings.HasPrefix(url, "youtube.com/watch?") {
      return greenTen
   }
   if unrated == 0 && good == 0 {
      return redTen
   }
   if unrated == 0 {
      return greenTen
   }
   if good == 0 {
      return redFive
   }
   return greenFive
}

func selectMb(tx winter.Tx, artist string) (string, error) {
   var mb string
   e := tx.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist,
   ).Scan(&mb)
   if e != nil {
      return "", e
   } else if mb == "" {
      return "", errors.New("mb_s missing")
   }
   return mb, nil
}

func localAlbum(tx winter.Tx, mb string) (map[string]winterLocal, error) {
   query, e := tx.Query(`
   select
      album_s,
      date_s,
      url_s,
      count(1) filter (where note_s = '') as unrated,
      count(1) filter (where note_s = 'good') as good
   from album_t
   natural join song_t
   natural join song_artist_t
   natural join artist_t
   where mb_s = ?
   group by album_n
   `, mb)
   if e != nil {
      return nil, e
   }
   var (
      locals = map[string]winterLocal{}
      q queryRow
   )
   for query.Next() {
      e = query.Scan(&q.album, &q.date, &q.url, &q.unrated, &q.good)
      if e != nil {
         return nil, e
      }
      locals[strings.ToUpper(q.album)] = winterLocal{
         color(q.url, q.unrated, q.good), q.date,
      }
   }
   return locals, nil
}

type queryRow struct {
   album string
   date string
   url string
   unrated int
   good int
}

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
