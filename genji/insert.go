package main

import (
   "encoding/json"
   "fmt"
   "github.com/genjidb/genji"
   "io/ioutil"
)

func Exec(o *genji.Tx, s string, a ...interface{}) error {
   fmt.Println(a)
   return o.Exec(s, a...)
}

func Insert(db *genji.DB) error {
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   // artist
   y, e := ioutil.ReadFile("artist.json")
   if e != nil {
      return e
   }
   artist_a := []Artist{}
   json.Unmarshal(y, &artist_a)
   e = Exec(tx, "CREATE TABLE artist_t")
   if e != nil {
      return e
   }
   for _, artist := range artist_a {
      e = Exec(tx, "INSERT INTO artist_t VALUES ?", artist)
      if e != nil {
         return e
      }
   }
   // album
   y, e = ioutil.ReadFile("album.json")
   if e != nil {
      return e
   }
   album_a := []Album{}
   json.Unmarshal(y, &album_a)
   e = Exec(tx, "CREATE TABLE album_t")
   if e != nil {
      return e
   }
   e = Exec(tx, "CREATE INDEX album_i ON album_t(id)")
   if e != nil {
      return e
   }
   for _, album := range album_a {
      e = Exec(tx, "INSERT INTO album_t VALUES ?", album)
      if e != nil {
         return e
      }
   }
   // song
   y, e = ioutil.ReadFile("song.json")
   if e != nil {
      return e
   }
   song_a := []Song{}
   json.Unmarshal(y, &song_a)
   e = Exec(tx, "CREATE TABLE song_t")
   if e != nil {
      return e
   }
   e = Exec(tx, "CREATE INDEX song_i ON song_t(id)")
   if e != nil {
      return e
   }
   for _, song := range song_a {
      e = Exec(tx, "INSERT INTO song_t VALUES ?", song)
      if e != nil {
         return e
      }
   }
   return tx.Commit()
}
