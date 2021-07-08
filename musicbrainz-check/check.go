package main

import (
   "database/sql"
   "errors"
   "fmt"
   "github.com/89z/mech/musicbrainz"
   "os"
   "strings"
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

func newLocalArtist(name, file string) (*localArtist, error) {
   db, err := sql.Open("sqlite3", file)
   if err != nil {
      return nil, err
   }
   defer db.Close()
   var artistID string
   if err := db.QueryRow(`
   SELECT mb_s FROM artist_t WHERE artist_s LIKE ?
   `, name).Scan(&artistID); err != nil {
      return nil, err
   } else if artistID == "" {
      return nil, errors.New("artistID missing")
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
      return nil, err
   }
   defer rows.Close()
   artist := &localArtist{
      artistID, make(map[string]localAlbum),
   }
   for rows.Next() {
      var a localAlbum
      err := rows.Scan(&a.title, &a.date, &a.url, &a.unrated, &a.good)
      if err != nil {
         return nil, err
      }
      artist.albums[a.date + a.title] = a
   }
   return artist, nil
}

func remoteAlbums(artistID string) ([]*musicbrainz.Release, error) {
   var (
      albums []*musicbrainz.Release
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

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release title.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type winterRemote struct {
   color, date, title string
   release map[string]bool
}

func color(url string, unrated, good int) string {
   const (
      block = "\u2587\u2587\u2587\u2587\u2587"
      gray = "\x1b[90m"
      green = "\x1b[92m"
      red = "\x1b[91m"
      reset = "\x1b[m"
   )
   if strings.HasPrefix(url, "youtube.com/watch?") {
      return green + block + block + reset
   }
   if unrated == 0 && good == 0 {
      return red + block + block + reset
   }
   if unrated == 0 {
      return green + block + block + reset
   }
   if good == 0 {
      return red + block + gray + block + reset
   }
   return green + block + gray + block + reset
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
   for n, group := range remotes {
      for release := range group.release {
         local, ok := locals[strings.ToUpper(release)]
         if ok {
            remotes[n].date = local.date
            remotes[n].color = local.color
         }
      }
   }
   sort.Slice(remotes, func(a, b int) bool {
      return remotes[a].date < remotes[b].date
   })
   for _, group := range remotes {
      fmt.Printf("%-10v | %10v | %v\n", group.date, group.color, group.title)
   }
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
