package main

import (
   "fmt"
   "time"
   "winter"
)

func selectAll(tx winter.Tx) error {
   then := time.Now().AddDate(-1, 0, 0)
   query, e := tx.Query(`
   select
      count(1) filter (where note_s = 'good') as good,
      artist_s
   from artist_t
   natural join song_artist_t
   natural join song_t
   where check_s < ?
   group by artist_n
   order by good
   `, then)
   if e != nil {
      return e
   }
   var (
      artist string
      count int
   )
   cmd, pipe, e := less()
   if e != nil {
      return e
   }
   defer cmd.Wait()
   defer pipe.Close()
   for query.Next() {
      e = query.Scan(&count, &artist)
      if e != nil {
         return e
      }
      fmt.Fprintln(pipe, count, "|", artist)
   }
   return nil
}
