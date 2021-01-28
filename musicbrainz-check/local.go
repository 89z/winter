package main

import (
   "database/sql"
   "github.com/89z/winter"
   "strings"
)

const (
   block = "\u2587\u2587\u2587\u2587\u2587"
   green_10 = "\x1b[92m" + block + block + "\x1b[m"
   green_5 = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[m"
   red_10 = "\x1b[91m" + block + block + "\x1b[m"
   red_5 = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[m"
)

func color(url_s string, unrated_n, good_n int) string {
   if winter.Pop(url_s) {
      return green_10
   }
   if unrated_n == 0 && good_n == 0 {
      return red_10
   }
   if unrated_n == 0 {
      return green_10
   }
   if good_n == 0 {
      return red_5
   }
   return green_5
}

func localAlbum(db *sql.DB, artist_s string) (map[string]local, error) {
   query, e := db.Query(`
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
   local_m := map[string]local{}
   for query.Next() {
      e = query.Scan(&album_s, &date_s, &url_s, &unrated_n, &good_n)
      if e != nil {
         return nil, e
      }
      local_m[strings.ToUpper(album_s)] = local{
         color(url_s, unrated_n, good_n), date_s,
      }
   }
   return local_m, nil
}
