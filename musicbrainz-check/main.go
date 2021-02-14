package main

import (
   "fmt"
   "log"
   "os"
   "sort"
   "strings"
   "winter"
)

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   tx, e := winter.NewTx(
      os.Getenv("WINTER"),
   )
   if e != nil {
      log.Fatal(e)
   }
   mb, e := selectMb(
      tx, os.Args[1],
   )
   if e != nil {
      log.Fatal(e)
   }
   // local albums
   locals, e := localAlbum(tx, mb)
   if e != nil {
      log.Fatal(e)
   }
   // remote albums
   remotes, e := remoteAlbum(mb)
   if e != nil {
      log.Fatal(e)
   }
   for n, group := range remotes {
      for release := range group.release {
         local, ok := locals[strings.ToUpper(release)]
         if ok {
            remotes[n].date = local.date
            remotes[n].color = local.color
         }
      }
   }
   sort.Slice(remotes, func(i, j int) bool {
      return remotes[i].date < remotes[j].date
   })
   for _, group := range remotes {
      fmt.Printf("%-10v | %10v | %v\n", group.date, group.color, group.title)
   }
}
