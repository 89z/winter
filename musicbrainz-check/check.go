package main

import (
   "database/sql"
   "errors"
   "fmt"
   "github.com/89z/mech/musicbrainz"
   "os"
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
   var artistID string
   if err := db.QueryRow(`
   SELECT mb_s FROM artist_t WHERE artist_s LIKE ?
   `, name).Scan(&artistID); err != nil {
      return localArtist{}, err
   } else if artistID == "" {
      return localArtist{}, errors.New("artistID missing")
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
   `, artistID)
   if err != nil {
      return localArtist{}, err
   }
   defer rows.Close()
   artist := localArtist{
      artistID, make(map[string]localAlbum),
   }
   for rows.Next() {
      var a localAlbum
      err := rows.Scan(&a.title, &a.date, &a.url, &a.unrated, &a.good)
      if err != nil {
         return localArtist{}, err
      }
      artist.albums[a.date + a.title] = a
   }
   return artist, nil
}

func remoteAlbums(artistID string) ([]musicbrainz.Release, error) {
   var (
      albums []musicbrainz.Release
      offset int
   )
   for {
      group, err := musicbrainz.GroupFromArtist(artistID, offset)
      if err != nil {
         return nil, err
      }
      for _, release := range group.Releases {
         if release.Date == "" {
            continue
         }
         if len(release.ReleaseGroup.SecondaryTypes) > 0 {
            continue
         }
         albums = append(albums, release)
      }
      offset += 100
      if offset >= group.ReleaseCount {
         break
      }
   }
   return albums, nil
}

func main() {
   if len(os.Args) != 2 {
      fmt.Println("musicbrainz-check <artist>")
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
   index, ok := remote[release.ReleaseGroup.ID]
   if ok {
      // add release to group
      remotes[index].release[release.Title] = true
   } else {
      // add group
      remotes = append(remotes, winterRemote{
         date: release.ReleaseGroup.FirstRelease,
         release: map[string]bool{release.Title: true},
         title: release.ReleaseGroup.Title,
      })
      remote[release.ReleaseGroup.ID] = len(remotes) - 1
   }
   */
}
