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
   if err != nil {
      return err
   }
   defer rows.Close()
   for rows.Next() {
      var (
         artist string
         count int
      )
      err := rows.Scan(&count, &artist)
      if err != nil {
         return err
      }
      fmt.Println(count, "|", artist)
   }
   return nil
}

func transact(tx *sql.Tx, name string, args []string) error {
   switch name {
   case "album":
      switch len(args) {
      case 1:
         return deleteAlbum(tx, args[0])
      case 2:
         return copyAlbum(tx, args[0], args[1])
      }
   case "artist":
      if len(args) == 0 {
         return selectAll(tx)
      }
      _, err := tx.Exec(`
      INSERT INTO artist_t (artist_s, check_s, mb_s) VALUES (?, '', '')
      `, args[0])
      return err
   case "check":
      _, err := tx.Exec(`
      UPDATE artist_t SET check_s = ? WHERE artist_n = ?
      `, args[1], args[0])
      return err
   case "date":
      _, err := tx.Exec(`
      UPDATE album_t SET date_s = ? WHERE album_n = ?
      `, args[1], args[0])
      return err
   case "mb":
      _, err := tx.Exec(`
      UPDATE artist_t SET mb_s = ? WHERE artist_n = ?
      `, args[1], args[0])
      return err
   case "note":
      _, err := tx.Exec(`
      UPDATE song_t SET note_s = ? WHERE song_n = ?
      `, args[1], args[0])
      return err
   case "url":
      _, err := tx.Exec(`
      UPDATE album_t SET url_s = ? WHERE album_n = ?
      `, args[1], args[0])
      return err
   default:
      return selectOne(tx, name)
   }
   return nil
}
