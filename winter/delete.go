package main

import (
   "database/sql"
   "winter/snow"
)

func DeleteAlbum(open_o *sql.DB, album_s string) error {
   query_o, e := open_o.Query(
      "select song_n from song_album_t where album_n = ?", album_s,
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
      e = snow.Delete(open_o, "song_album_t where song_n = ?", song_n)
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
