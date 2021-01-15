package main

import (
   "database/sql"
   "fmt"
   "time"
)

func SelectAll(tx *sql.Tx) error {
   then := time.Now().AddDate(-1, 0, 0)
   query_o, e := tx.Query(`
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
   less, pipe, e := Less()
   if e != nil {
      return e
   }
   defer less.Wait()
   defer pipe.Close()
   for query_o.Next() {
      e = query_o.Scan(&count_n, &artist_s)
      if e != nil {
         return e
      }
      fmt.Fprintln(pipe, count_n, "|", artist_s)
   }
   return nil
}
