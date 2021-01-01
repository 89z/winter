package main

import (
   "dbstorm"
   "encoding/json"
   "github.com/asdine/storm/v3"
   "log"
   "os"
)

func Encode(v interface{}, s string) error {
   create_o, e := os.Create(s)
   if e != nil {
      return e
   }
   json_o := json.NewEncoder(create_o)
   json_o.SetIndent("", " ")
   return json_o.Encode(v)
}

func main() {
   db, e := storm.Open("thunder.db")
   if e != nil {
      log.Fatal(e)
   }
   // artist
   artist_a := []dbstorm.Artist{}
   e = db.All(&artist_a)
   if e != nil {
      log.Fatal(e)
   }
   e = Encode(artist_a, "artist.json")
   if e != nil {
      log.Fatal(e)
   }
   // album
   album_a := []dbstorm.Album{}
   e = db.All(&album_a)
   if e != nil {
      log.Fatal(e)
   }
   e = Encode(album_a, "album.json")
   if e != nil {
      log.Fatal(e)
   }
   // song
   song_a := []dbstorm.Song{}
   e = db.All(&song_a)
   if e != nil {
      log.Fatal(e)
   }
   e = Encode(song_a, "song.json")
   if e != nil {
      log.Fatal(e)
   }
}
