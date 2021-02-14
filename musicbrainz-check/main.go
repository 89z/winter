package main

import (
   "errors"
   "log"
   "os"
   "winter"
)

func getLocal(name, file string) (localArtist, error) {
   tx, e := winter.NewTx(file)
   if e != nil {
      return localArtist{}, e
   }
   var artistId string
   e := tx.QueryRow(
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

type localArtist struct {
   artistId string
   albums []localAlbum
}

type localAlbum struct {
   date string
   good int
   title string
   unrated int
   url string
}

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   name, file := os.Args[1], os.Getenv("WINTER")
   locals, e := getLocal(name, file)
   if e != nil {
      log.Fatal(e)
   }
   locals := map[string]winterLocal{}
   for query.Next() {
      var q queryRow
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
