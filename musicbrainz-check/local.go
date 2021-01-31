package main

import (
   "database/sql"
   "github.com/89z/winter"
   "strings"
)

const (
   block = "\u2587\u2587\u2587\u2587\u2587"
   greenFive = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[m"
   greenTen = "\x1b[92m" + block + block + "\x1b[m"
   redFive = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[m"
   redTen = "\x1b[91m" + block + block + "\x1b[m"
)

func color(url string, unrated, good int) string {
   if winter.Pop(url) {
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
      count(1) filter (where note_s = '') as unrated_n,
      count(1) filter (where note_s = 'good') as good_n
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
      album_s string
      date_s string
      good_n int
      unrated_n int
      url_s string
   )
   local_m := map[string]winterLocal{}
   for query.Next() {
      e = query.Scan(&album_s, &date_s, &url_s, &unrated_n, &good_n)
      if e != nil {
         return nil, e
      }
      local_m[strings.ToUpper(album_s)] = winterLocal{
         color(url_s, unrated_n, good_n), date_s,
      }
   }
   return local_m, nil
}
