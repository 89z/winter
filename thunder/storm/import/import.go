package main

import (
   "encoding/json"
   "github.com/asdine/storm/v3"
   "os"
   "thunder/clap"
)

var (
   e error
   o *os.File
)

func Import() error {
   db, e := storm.Open("thunder.db")
   if e != nil {
      return e
   }
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   // decode song
   o, e = os.Open("song.json")
   if e != nil {
      return e
   }
   song_a := clap.Slice{}
   json.NewDecoder(o).Decode(&song_a)
   // save songs
   for n := range song_a {
      song_m := song_a.M(n)
      song := clap.Song{
         Name: song_m.S("song_s"), Note: song_m.S("note_s"),
      }
      e = clap.Save(tx, &song)
      if e != nil {
         return e
      }
      old_song := clap.OldSong{
         Old: int(song_m.N("song_n")), Album: int(song_m.N("album_n")),
      }
      e = clap.Save(tx, &old_song)
      if e != nil {
         return e
      }
   }
   // decode song artist
   o, e = os.Open("song_artist.json")
   if e != nil {
      return e
   }
   song_artist_a := clap.Slice{}
   json.NewDecoder(o).Decode(&song_artist_a)
   for n := range song_artist_a {
      song_artist_m := song_artist_a.M(n)
      song_artist := clap.SongArtist{
         Song: int(song_artist_m.N("song_n")),
         Artist: int(song_artist_m.N("artist_n")),
      }
      e = clap.Save(tx, &song_artist)
      if e != nil {
         return e
      }
   }
   // decode album
   o, e = os.Open("album.json")
   if e != nil {
      return e
   }
   album_a := clap.Slice{}
   json.NewDecoder(o).Decode(&album_a)
   for n := range album_a {
      album_m := album_a.M(n)
      album := clap.Album{
         Date: album_m.S("date_s"),
         Name: album_m.S("album_s"),
         URL: album_m.S("url_s"),
      }
      e = clap.Save(tx, &album)
      if e != nil {
         return e
      }
      old_album := clap.OldAlbum{
         Old: int(album_m.N("album_n")),
      }
      e = clap.Save(tx, &old_album)
      if e != nil {
         return e
      }
   }
   // decode artist
   o, e = os.Open("artist.json")
   if e != nil {
      return e
   }
   artist_a := clap.Slice{}
   json.NewDecoder(o).Decode(&artist_a)
   for n := range artist_a {
      artist_m := artist_a.M(n)
      artist := clap.Artist{
         Check: artist_m.S("check_s"),
         MB: artist_m.S("mb_s"),
         Name: artist_m.S("artist_s"),
      }
      e = clap.Save(tx, &artist)
      if e != nil {
         return e
      }
   }
   return tx.Commit()
}
