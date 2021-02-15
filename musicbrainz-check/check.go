package main

import (
   "errors"
   "net/url"
   "winter"
)


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
   var artist localArtist
   for query.Next() {
      var alb localAlbum
      e = query.Scan(&alb.title, &alb.date, &alb.url, &alb.unrated, &alb.good)
      if e != nil {
         return localArtist{}, e
      }
      artist.albums = append(artist.albums, alb)
   }
   return artist, nil
}


func newRemoteArtist(artistId string) (remoteArtist, error) {
   value := url.Values{}
   value.Set("fmt", "json")
   value.Set("inc", "release-groups")
   value.Set("limit", "100")
   value.Set("status", "official")
   value.Set("type", "album")
   value.Set("artist", artistId)
   for {
      get, e := http.Get(
         "http://musicbrainz.org/ws/2/release?" + value.Encode(),
      )
      if e != nil {
         return remoteArtist{}, e
      }
      var album remoteAlbum
      e = json.NewDecoder(get.Body).Decode(&album)
      if e != nil {
         return remoteArtist{}, e
      }
      for _, release := range album.Releases {
         if release.Date == "" {
            continue
         }
         if len(release.Group.SecondaryTypes) > 0 {
            continue
         }
      }
   }
}

type localAlbum struct {
   date string
   good int
   title string
   unrated int
   url string
}

type localArtist struct {
   artistId string
   albums []localAlbum
}

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
