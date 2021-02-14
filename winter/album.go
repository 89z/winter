package main

import (
   "database/sql"
   "winter"
)

func copyAlbum(tx *sql.Tx, source , dest string) error {
   var note, song, url_s string
   // COPY URL
   e := tx.QueryRow(
      "select url_s from album_t where album_n = ?", source,
   ).Scan(&url_s)
   if e != nil {
      return e
   }
   // PASTE URL
   _, e = winter.Exec(
      tx, "update album_t set url_s = ? where album_n = ?", url_s, dest,
   )
   if e != nil {
      return e
   }
   // COPY NOTES
   query, e := tx.Query(
      "select song_s, note_s from song_t where album_n = ?", source,
   )
   if e != nil {
      return e
   }
   songs := map[string]string{}
   for query.Next() {
      e = query.Scan(&song, &note)
      if e != nil {
         return e
      }
      songs[song] = note
   }
   // PASTE NOTES
   for song, note := range songs {
      _, e = winter.Exec(tx, `
      update song_t set note_s = ?
      where album_n = ? and song_s = ? COLLATE NOCASE
      `, note, dest, song)
      if e != nil {
         return e
      }
   }
   return nil
}

func deleteAlbum(tx *sql.Tx, album string) error {
   query, e := tx.Query("select song_n from song_t where album_n = ?", album)
   if e != nil {
      return e
   }
   var (
      song int
      songs []int
   )
   for query.Next() {
      e = query.Scan(&song)
      if e != nil {
         return e
      }
      songs = append(songs, song)
   }
   for _, song := range songs {
      e = winter.Delete(tx, "song_t where song_n = ?", song)
      if e != nil {
         return e
      }
      e = winter.Delete(tx, "song_artist_t where song_n = ?", song)
      if e != nil {
         return e
      }
   }
   e = winter.Delete(tx, "album_t where album_n = ?", album)
   if e != nil {
      return e
   }
   return nil
}
