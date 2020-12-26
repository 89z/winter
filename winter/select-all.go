package main

import (
   "database/sql"
   "time"
)

func SelectAll(open_o *sql.DB) error {
   then := time.Now().AddDate(-1, 0, 0)
   query_o, e := open_o.Query(`
   select
      count(1) filter (where note_s = 'good') as count_n,
      artist_s
   from artist_t
   natural join song_artist_t
   natural join song_t
   where check_s < ?
   group by artist_n
   order by count_n
   `, then)
   if e != nil {
      return e
   }
   var (
      artist_s string
      count_n int
   )
   for query_o.Next() {
      e = query_o.Scan(&count_n, &artist_s)
      if e != nil {
         return e
      }
      println(count_n, "|", artist_s)
   }
   return nil
}
