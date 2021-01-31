package main

import (
   "database/sql"
   "fmt"
   "github.com/89z/x"
   "log"
   "os"
   "sort"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

var mb string

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist := os.Args[1]
   db, e := sql.Open(
      "sqlite3", os.Getenv("WINTER"),
   )
   x.Check(e)
   e = db.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist,
   ).Scan(&mb)
   x.Check(e)
   if mb == "" {
      log.Fatal("mb_s missing")
   }
   // local albums
   locals, e := localAlbum(db, mb)
   x.Check(e)
   // remote albums
   remote, e := remoteAlbum(mb)
   x.Check(e)
   for n, group := range remote {
      for release := range group.release {
         local, b := locals[strings.ToUpper(release)]
         if b {
            remote[n].date = local.date
            remote[n].color = local.color
         }
      }
   }
   sort.Slice(remote, func(n, n2 int) bool {
      return remote[n].date < remote[n2].date
   })
   for _, group := range remote {
      fmt.Printf("%-10v | %10v | %v\n", group.date, group.color, group.title)
   }
}
