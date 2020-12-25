package main

import (
   "database/sql"
   "winter/snow"
)

const (
   U2587 = "\u2587\u2587\u2587\u2587\u2587"
   GREEN_10 = "\x1b[92m" + U2587 + U2587 + "\x1b[m"
   GREEN_5 = "\x1b[92m" + U2587 + "\x1b[90m" + U2587 + "\x1b[m"
   RED_10 = "\x1b[91m" + U2587 + U2587 + "\x1b[m"
   RED_5 = "\x1b[91m" + U2587 + "\x1b[90m" + U2587 + "\x1b[m"
)

func LocalAlbum(open_o *sql.DB, artist_s string) (map[string]string, error) {
   query_o, e := open_o.Query(`
   select
      album_s,
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
      url_s string
      unrated_n int
      good_n int
   )
   local_m := map[string]string{}
   for query_o.Next() {
      e = query_o.Scan(&album_s, &url_s, &unrated_n, &good_n)
      if e != nil {
         return nil, e
      }
      if snow.Pop(url_s) {
         local_m[album_s] = GREEN_10
         continue
      }
      if unrated_n == 0 && good_n == 0 {
         local_m[album_s] = RED_10
         continue
      }
      if unrated_n == 0 {
         local_m[album_s] = GREEN_10
         continue
      }
      if good_n == 0 {
         local_m[album_s] = RED_5
         continue
      }
      local_m[album_s] = GREEN_5
   }
   return local_m, nil
}
