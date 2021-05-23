package main

import (
   "database/sql"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "os"
   "strconv"
   _ "github.com/mattn/go-sqlite3"
)

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

func newLocalArtist(name, file string) (localArtist, error) {
   db, err := sql.Open("sqlite3", file)
   if err != nil {
      return localArtist{}, err
   }
   defer db.Close()
   var artistId string
   if err := db.QueryRow(`
   SELECT mb_s FROM artist_t WHERE artist_s LIKE ?
   `, name).Scan(&artistId); err != nil {
      return localArtist{}, err
   } else if artistId == "" {
      return localArtist{}, errors.New("artistId missing")
   }
   rows, err := db.Query(`
   SELECT
      album_s,
      date_s,
      url_s,
      count(1) filter (where note_s = '') as unrated,
      count(1) filter (where note_s = 'good') as good
   FROM album_t
   NATURAL JOIN song_t
   NATURAL JOIN song_artist_t
   NATURAL JOIN artist_t
   WHERE mb_s = ?
   GROUP BY album_n
   `, artistId)
   if err != nil {
      return localArtist{}, err
   }
   defer rows.Close()
   artist := localArtist{
      artistId, make(map[string]localAlbum),
   }
   for rows.Next() {
      var alb localAlbum
      err := rows.Scan(&alb.title, &alb.date, &alb.url, &alb.unrated, &alb.good)
      if err != nil {
         return localArtist{}, err
      }
      artist.albums[alb.date + alb.title] = alb
   }
   return artist, nil
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
      json.NewDecoder(res.Body).Decode(&artist)
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

type remoteArtist struct {
   ReleaseCount int `json:"release-count"`
   Releases []remoteAlbum
}

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      return
   }
   name, file := os.Args[1], os.Getenv("WINTER")
   local, err := newLocalArtist(name, file)
   if err != nil {
      panic(err)
   }
   remote, err := remoteAlbums(local.id)
   if err != nil {
      panic(err)
   }
   fmt.Println(remote)
   /*
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
   */
}
