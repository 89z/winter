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
   local := map[string]winterLocal{}
   for query.Next() {
      var q queryRow
      e = query.Scan(&q.album, &q.date, &q.url, &q.unrated, &q.good)
      if e != nil {
         return nil, e
      }
      local[strings.ToUpper(q.album)] = winterLocal{
         color(q.url, q.unrated, q.good), q.date,
      }
   }
   return local, nil
   */
}
