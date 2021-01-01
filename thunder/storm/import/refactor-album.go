package main

import (
   "fmt"
   "github.com/asdine/storm/v3"
   "thunder/clap"
)

func RefactorAlbum() error {
   db, e := storm.Open("thunder.db")
   if e != nil {
      return e
   }
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   // get old songs
   old_song_a := []clap.OldSong{}
   e = tx.All(&old_song_a)
   if e != nil {
      return e
   }
   for _, old_song := range old_song_a {
      // get old album
      old_album := clap.OldAlbum{}
      e = tx.One("Old", old_song.Album, &old_album)
      if e != nil {
         return fmt.Errorf("old_song %v %v", old_song.Album, e)
      }
      // get new album
      album := clap.Album{}
      e = tx.One("ID", old_album.ID, &album)
      if e != nil {
         return fmt.Errorf("old_album %v %v", old_album.ID, e)
      }
      // append
      album.Song = append(album.Song, old_song.ID)
      // update
      e = clap.Update(tx, &album)
      if e != nil {
         return e
      }
   }
   return tx.Commit()
}
