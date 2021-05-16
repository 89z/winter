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
   req, err := http.NewRequest("GET", "http://musicbrainz.org/ws/2/release", nil)
   if err != nil { return nil, err }
   val := req.URL.Query()
   val.Set("fmt", "json")
   val.Set("inc", "release-groups")
   val.Set("limit", "100")
   val.Set("status", "official")
   val.Set("type", "album")
   val.Set("artist", artistId)
   for {
      req.URL.RawQuery = val.Encode()
      res, err := new(http.Client).Do(req)
      if err != nil { return nil, err }
      var artist remoteArtist
      err = json.NewDecoder(res.Body).Decode(&artist)
      if err != nil { return nil, err }
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
   tx, err := winter.NewTx(file)
   if err != nil {
      return localArtist{}, err
   }
   var artistId string
   err = tx.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", name,
   ).Scan(&artistId)
   if err != nil {
      return localArtist{}, err
   } else if artistId == "" {
      return localArtist{}, errors.New("artistId missing")
   }
   query, err := tx.Query(`
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
   if err != nil {
      return localArtist{}, err
   }
   artist := localArtist{
      artistId, map[string]localAlbum{},
   }
   for query.Next() {
      var alb localAlbum
      err = query.Scan(&alb.title, &alb.date, &alb.url, &alb.unrated, &alb.good)
      if err != nil {
         return localArtist{}, err
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
