package main

import (
   "database/sql"
   "strings"
   "winter"
)

const (
   BLOCK = "\u2587\u2587\u2587\u2587\u2587"
   GREEN_10 = "\x1b[92m" + BLOCK + BLOCK + "\x1b[m"
   GREEN_5 = "\x1b[92m" + BLOCK + "\x1b[90m" + BLOCK + "\x1b[m"
   RED_10 = "\x1b[91m" + BLOCK + BLOCK + "\x1b[m"
   RED_5 = "\x1b[91m" + BLOCK + "\x1b[90m" + BLOCK + "\x1b[m"
)

func Color(url_s string, unrated_n, good_n int) string {
   if winter.Pop(url_s) {
      return GREEN_10
   }
   if unrated_n == 0 && good_n == 0 {
      return RED_10
   }
   if unrated_n == 0 {
      return GREEN_10
   }
   if good_n == 0 {
      return RED_5
   }
   return GREEN_5
}

func LocalAlbum(db *sql.DB, artist_s string) (map[string]Local, error) {
   query_o, e := db.Query(`
   select
      album_s,
      date_s,
      url_s,
      count(1) filter (where note_s = '') as unrated_n,
      count(1) filter (where note_s = 'good') as good_n
   from album_t
   natural join song_t
   natural join song_artist_t
   natural join artist_t
   where artist_s LIKE ?
   group by album_n
   `, artist_s)
   if e != nil {
      return nil, e
   }
   var (
      album_s string
      date_s string
      good_n int
      unrated_n int
      url_s string
   )
   local_m := map[string]Local{}
   for query_o.Next() {
      e = query_o.Scan(&album_s, &date_s, &url_s, &unrated_n, &good_n)
      if e != nil {
         return nil, e
      }
      local_m[strings.ToUpper(album_s)] = Local{
         Color(url_s, unrated_n, good_n), date_s,
      }
   }
   return local_m, nil
}

type Local struct {
   Color string
   Date string
}
