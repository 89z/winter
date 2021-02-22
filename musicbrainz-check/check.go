package main

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "winter"
)

type remoteAlbum struct {
   Date string
   Group struct {
      FirstRelease string `json:"first-release-date"`
      Id string
      SecondaryTypes []string `json:"secondary-types"`
      Title string
   } `json:"release-group"`
   Title string
}

type remoteArtist struct {
   ReleaseCount int `json:"release-count"`
   Releases []remoteAlbum
}

func remoteAlbums(artistId string) ([]remoteAlbum, error) {
   value := url.Values{}
   value.Set("fmt", "json")
   value.Set("inc", "release-groups")
   value.Set("limit", "100")
   value.Set("status", "official")
   value.Set("type", "album")
   value.Set("artist", artistId)
   var (
      albums []remoteAlbum
      offset int
   )
   for {
      get, e := http.Get(
         "http://musicbrainz.org/ws/2/release?" + value.Encode(),
      )
      if e != nil {
         return nil, e
      }
      var artist remoteArtist
      e = json.NewDecoder(get.Body).Decode(&artist)
      if e != nil {
         return nil, e
      }
      for _, release := range artist.Releases {
         if release.Date == "" {
            continue
         }
         if len(release.Group.SecondaryTypes) > 0 {
            continue
         }
         albums = append(albums, release)
      }
      offset += 100
      if offset >= artist.ReleaseCount {
         break
      }
      value.Set(
         "offset", fmt.Sprint(offset),
      )
   }
   return albums, nil
}

func newLocalArtist(name, file string) (localArtist, error) {
   tx, e := winter.NewTx(file)
   if e != nil {
      return localArtist{}, e
   }
   var artistId string
   e = tx.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", name,
   ).Scan(&artistId)
   if e != nil {
      return localArtist{}, e
   } else if artistId == "" {
      return localArtist{}, errors.New("artistId missing")
   }
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
   `, artistId)
   if e != nil {
      return localArtist{}, e
   }
   artist := localArtist{
      artistId,
      map[string]localAlbum{},
   }
   for query.Next() {
      var alb localAlbum
      e = query.Scan(&alb.title, &alb.date, &alb.url, &alb.unrated, &alb.good)
      if e != nil {
         return localArtist{}, e
      }
      artist.albums[alb.date + alb.title] = alb
   }
   return artist, nil
}

type localAlbum struct {
   date string
   good int
   title string
   unrated int
   url string
}

type localArtist struct {
   id string
   albums map[string]localAlbum
}
