package main

import (
   "encoding/json"
   "errors"
   "net/http"
   "strconv"
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
   var (
      albums []remoteAlbum
      offset int
   )
   req, e := http.NewRequest("GET", "http://musicbrainz.org/ws/2/release", nil)
   if e != nil { return nil, e }
   val := req.URL.Query()
   val.Set("fmt", "json")
   val.Set("inc", "release-groups")
   val.Set("limit", "100")
   val.Set("status", "official")
   val.Set("type", "album")
   val.Set("artist", artistId)
   for {
      req.URL.RawQuery = val.Encode()
      res, e := http.DefaultClient.Do(req)
      if e != nil { return nil, e }
      var artist remoteArtist
      e = json.NewDecoder(res.Body).Decode(&artist)
      if e != nil { return nil, e }
      for _, release := range artist.Releases {
         if release.Date == "" { continue }
         if len(release.Group.SecondaryTypes) > 0 { continue }
         albums = append(albums, release)
      }
      offset += 100
      if offset >= artist.ReleaseCount { break }
      val.Set("offset", strconv.Itoa(offset))
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
      artistId, map[string]localAlbum{},
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
   good, unrated int
   date, title, url string
}

type localArtist struct {
   id string
   albums map[string]localAlbum
}
