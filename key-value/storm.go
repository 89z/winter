package main

import (
   "github.com/asdine/storm/v3"
   "log"
)

func main() {
   db, e := storm.Open("winter.db")
   if e != nil {
      log.Fatal(e)
   }
   defer db.Close()
   e = Save(db)
   if e != nil {
      log.Fatal(e)
   }
   // QUERY
}
