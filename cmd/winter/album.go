package main

import (
   "database/sql"
   "winter"
)

func CopyAlbum(db *sql.DB, source , dest string) error {
   var note_s, song_s, url_s string
   tx, e := db.Begin()
   if e != nil {
      return e
   }
   // COPY URL
   e = tx.QueryRow(
      "select url_s from album_t where album_n = ?", source,
   ).Scan(&url_s)
   if e != nil {
      return e
   }
   // PASTE URL
   e = winter.Update(
      tx, "album_t set url_s = ? where album_n = ?", url_s, dest,
   )
   if e != nil {
      return e
   }
   // COPY NOTES
   query_o, e := tx.Query(
      "select song_s, note_s from song_t where album_n = ?", source,
   )
   if e != nil {
      return e
   }
   song_m := map[string]string{}
   for query_o.Next() {
      e = query_o.Scan(&song_s, &note_s)
      if e != nil {
         return e
      }
      song_m[song_s] = note_s
   }
   // PASTE NOTES
   for song_s, note_s := range song_m {
      e = winter.Update(
         tx,
         "song_t set note_s = ? where album_n = ? and song_s = ? COLLATE NOCASE",
         note_s,
         dest,
         song_s,
      )
      if e != nil {
         return e
      }
   }
   return tx.Commit()
}

func DeleteAlbum(db *sql.DB, album_s string) error {
   tx, e := db.Begin()
   if e != nil {
      return e
   }
   query_o, e := tx.Query(
      "select song_n from song_t where album_n = ?", album_s,
   )
   if e != nil {
      return e
   }
   song_n, song_a := 0, []int{}
   for query_o.Next() {
      e = query_o.Scan(&song_n)
      if e != nil {
         return e
      }
      song_a = append(song_a, song_n)
   }
   for _, song_n := range song_a {
      e = winter.Delete(tx, "song_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
      e = winter.Delete(tx, "song_artist_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
   }
   e = winter.Delete(tx, "album_t where album_n = ?", album_s)
   if e != nil {
      return e
   }
   return tx.Commit()
}
