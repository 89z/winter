package main

import (
   "github.com/genjidb/genji"
   "log"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      println("thunder <insert | select>")
      os.Exit(1)
   }
   db, e := genji.Open("thunder.db")
   if e != nil {
      log.Fatal(e)
   }
   defer db.Close()
   if os.Args[1] == "insert" {
      e = Insert(db)
   } else {
      e = Select(db)
   }
   if e != nil {
      log.Fatal(e)
   }
}
