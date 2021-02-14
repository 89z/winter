package main

import (
   "os"
   "log"
   "winter"
)

func localAlbums(artist, file string) error {
   tx, e := winter.NewTx(file)
   if e != nil {
      return e
   }
   var mb string
   e := tx.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist,
   ).Scan(&mb)
   if e != nil {
      return e
   } else if mb == "" {
      return errors.New("mb_s missing")
   }
   // FIXME
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

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist, file := os.Args[1], os.Getenv("WINTER")
   locals, e := localAlbums(artist, file)
   if e != nil {
      log.Fatal(e)
   }
}
