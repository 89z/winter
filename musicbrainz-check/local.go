package main

import (
   "database/sql"
   "strings"
)

func color(url string, unrated, good int) string {
   const (
      block = "\u2587\u2587\u2587\u2587\u2587"
      greenFive = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[m"
      greenTen = "\x1b[92m" + block + block + "\x1b[m"
      redFive = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[m"
      redTen = "\x1b[91m" + block + block + "\x1b[m"
   )
   if strings.HasPrefix(url, "youtube.com/watch?") {
      return greenTen
   }
   if unrated == 0 && good == 0 {
      return redTen
   }
   if unrated == 0 {
      return greenTen
   }
   if good == 0 {
      return redFive
   }
   return greenFive
}

func localAlbum(db *sql.DB, mb string) (map[string]winterLocal, error) {
   query, e := db.Query(`
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

type queryRow struct {
   album string
   date string
   url string
   unrated int
   good int
}
