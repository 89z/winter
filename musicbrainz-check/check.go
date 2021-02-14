package main

import (
   "errors"
   "winter"
)

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

