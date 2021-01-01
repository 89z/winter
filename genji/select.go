package main

import (
   "fmt"
   "github.com/genjidb/genji"
   "github.com/genjidb/genji/document"
)

func Select(db *genji.DB) error {
   tx, e := db.Begin(true)
   if e != nil {
      return e
   }
   // artist
   d, e := tx.QueryDocument("SELECT album from artist_t where id = 124")
   if e != nil {
      return e
   }
   album_a := []int{}
   e = document.Scan(d, &album_a)
   if e != nil {
      return e
   }
   for _, album_n := range album_a {
      // album
      d, e = tx.QueryDocument("SELECT song from album_t where id = ?", album_n)
      if e != nil {
         return e
      }
      song_a := []int{}
      e = document.Scan(d, &song_a)
      if e != nil {
         return e
      }
      for _, song_n := range song_a {
         // song
         d, e = tx.QueryDocument("SELECT name from song_t where id = ?", song_n)
         if e != nil {
            return e
         }
         var name string
         e = document.Scan(d, &name)
         if e != nil {
            return e
         }
         fmt.Println(name)
      }
   }
   return tx.Commit()
}
