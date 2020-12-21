package main

import (
   "database/sql"
   "winter/snow"
)

func DeleteAlbum(open_o *sql.DB, album_s string) error {
   e := snow.Delete("album_t where album_n = ?", album_s)
   if e != nil {
      return e
   }
   query_o, e := open_o.Query(
      "SELECT song_n FROM song_album_t WHERE album_n = ?", album_s,
   )
   if e != nil {
      return e
   }
   var song_n int
   for query_o.Next() {
      e = query_o.Scan(&song_n)
      if e != nil {
         return e
      }
      e = snow.Delete("song_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
      e = snow.Delete("song_album_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
      e = snow.Delete("song_artist_t where song_n = ?", song_n)
      if e != nil {
         return e
      }
   }
}
