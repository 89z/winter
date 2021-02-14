package main

import (
   "log"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   name, file := os.Args[1], os.Getenv("WINTER")
   local, e := newLocalArtist(name, file)
   if e != nil {
      log.Fatal(e)
   }
   remote, e := newRemoteArtist(local.artistId)
   if e != nil {
      log.Fatal(e)
   }
   /*
   local[strings.ToUpper(q.album)] = winterLocal{
      color(q.url, q.unrated, q.good), q.date,
   }
   */
}
