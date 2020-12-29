package main

import (
   "database/sql"
   "winter/snow"
)

func CopyAlbum(open_o *sql.DB, source , dest string) error {
   var note_s, song_s, url_s string
   // COPY URL
   e := open_o.QueryRow(
      "select url_s from album_t where album_n = ?", source,
   ).Scan(&url_s)
   if e != nil {
      return e
   }
   // PASTE URL
   e = snow.Update(
      open_o, "album_t set url_s = ? where album_n = ?", url_s, dest,
   )
   if e != nil {
      return e
   }
   // COPY NOTES
   query_o, e := open_o.Query(
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
      e = snow.Update(
         open_o,
         "song_t set note_s = ? where album_n = ? and song_s = ? COLLATE NOCASE",
         note_s,
         dest,
         song_s,
      )
      if e != nil {
         return e
      }
   }
   return nil
}

func DeleteAlbum(open_o *sql.DB, album_s string) error {
   query_o, e := open_o.Query(
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
      e = snow.Delete(open_o, "song_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
      e = snow.Delete(open_o, "song_artist_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
   }
   e = snow.Delete(open_o, "album_t where album_n = ?", album_s)
   if e != nil {
      return e
   }
   return nil
}
