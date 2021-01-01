package main

import (
   "log"
   "os"
)

func main() {
   if len(os.Args) == 1 {
      println(`json import
json refactor artist
json refactor album`)
      os.Exit(1)
   }
   var (
      e error
      key_s = os.Args[1]
   )
   if key_s == "import" {
      e = Import()
   } else if os.Args[2] == "artist" {
      e = RefactorArtist()
   } else {
      e = RefactorAlbum()
   }
   if e != nil {
      log.Fatal(e)
   }
}
