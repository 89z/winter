package main

import (
   "fmt"
   "github.com/asdine/storm/v3"
   "thunder/clap"
)

func RefactorArtist() error {
   db, e := storm.Open("thunder.db")
   if e != nil {
      return e
   }
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   album_a := []clap.Album{}
   e = tx.All(&album_a)
   if e != nil {
      return e
   }
   for _, album := range album_a {
      // get old album
      old_album := clap.OldAlbum{}
      e = tx.One("ID", album.ID, &old_album)
      if e != nil {
         return fmt.Errorf("album %v %v", album.ID, e)
      }
      // get old song
      old_song := clap.OldSong{}
      e = tx.One("Album", old_album.Old, &old_song)
      if e != nil {
         return fmt.Errorf("old_album %v %v", old_album.Old, e)
      }
      // get song artists
      song_artist_a := []clap.SongArtist{}
      e = tx.Find("Song", old_song.Old, &song_artist_a)
      if e != nil {
         return fmt.Errorf("old_song %v %v", old_song.Old, e)
      }
      for _, song_artist := range song_artist_a {
         // get artist
         artist := clap.Artist{}
         e = tx.One("ID", song_artist.Artist, &artist)
         if e != nil {
            return fmt.Errorf("song_artist %v %v", song_artist.Artist, e)
         }
         // append
         artist.Album = append(artist.Album, album.ID)
         // update
         e = clap.Update(tx, &artist)
         if e != nil {
            return e
         }
      }
   }
   return tx.Commit()
}
