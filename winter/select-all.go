package main

import (
   "fmt"
   "time"
   "winter"
)

func selectAll(tx winter.Tx) error {
   then := time.Now().AddDate(-1, 0, 0)
   query, err := tx.Query(`
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
   if err != nil { return err }
   var (
      artist string
      count int
   )
   for query.Next() {
      err = query.Scan(&count, &artist)
      if err != nil { return err }
      fmt.Println(count, "|", artist)
   }
   return nil
}
