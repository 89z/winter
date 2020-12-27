package main

import (
   "github.com/asdine/storm/v3"
   "log"
)

type Artist struct {
   ID int `storm:"increment"`
   Name string
   Song []int
}

type Song struct {
   ID int `storm:"increment"`
   Name string
   Artist []int
}

func main() {
   db, e := storm.Open("winter.db")
   defer db.Close()
   var (
      artist Artist
      song Song
   )
   // 1
   artist = Artist{
      Name: "Zero 7",
      Song: []int{0, 1},
   }
   e = db.Save(&artist)
   if e != nil {
      log.Fatal(e)
   }
   // 2
   artist = Artist{
      Name: "Sia",
      Song: []int{0, 2},
   }
   e = db.Save(&artist)
   if e != nil {
      log.Fatal(e)
   }
   // 3
   song = Song{
      Name: "Destiny",
      Artist: []int{0, 1}
   }
   e = db.Save(&song)
   if e != nil {
      log.Fatal(e)
   }
   // 4
   song = Song{
      Name: "Give It Away",
      Artist: []int{0},
   }
   e = db.Save(&song)
   if e != nil {
      log.Fatal(e)
   }
}
