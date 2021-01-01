package dbstorm

import (
   "fmt"
   "github.com/asdine/storm/v3"
   "os"
)

func SelectOne(db *storm.DB, artist_s string) error {
   artist_o := Artist{}
   e := db.One("Name", artist_s, &artist_o)
   if e != nil {
      return e
   }
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   pipe := os.Stdout
   for _, album_n := range artist_o.Album {
      album_o := Album{}
      e = tx.One("ID", album_n, &album_o)
      if e != nil {
         return e
      }
      for _, song_n := range album_o.Song {
         song_o := Song{}
         e = tx.One("ID", song_n, &song_o)
         if e != nil {
            return e
         }
         fmt.Fprintln(pipe, song_o.ID)
      }
   }
   return tx.Commit()
}
