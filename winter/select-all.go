package main

import (
   "database/sql"
   "fmt"
   "time"
)

func selectAll(tx *sql.Tx) error {
   then := time.Now().AddDate(-1, 0, 0)
   rows, err := tx.Query(`
   SELECT
      count(1) filter (WHERE note_s = 'good') AS good,
      artist_s
   FROM artist_t
   NATURAL JOIN song_artist_t
   NATURAL JOIN song_t
   WHERE check_s < ?
   GROUP BY artist_n
   ORDER BY good
   `, then)
   if err != nil { return err }
   defer rows.Close()
   for rows.Next() {
      var (
         artist string
         count int
      )
      err := rows.Scan(&count, &artist)
      if err != nil { return err }
      fmt.Println(count, "|", artist)
   }
   return nil
}
