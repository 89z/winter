package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
   "sort"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

/* Regarding the title and date:

For the title, we will display the remote Group title, but we also need to get
the remote Release titles to match against the local Release titles.

For the date, if we have a local match, use that date. Otherwise, use use the
remote Group date */
type winterRemote struct{
   color string
   date string
   release map[string]bool
   title string
}

type winterLocal struct{
   color string
   date string
}

type mbRelease struct{
   ReleaseCount int `json:"release-count"`
   Releases []struct{
      Date string
      Group struct{
         FirstRelease string `json:"first-release-date"`
         Id string
         SecondaryTypes []string `json:"secondary-types"`
         Title string
      } `json:"release-group"`
      Title string
   }
}

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist := os.Args[1]
   db, e := sql.Open(
      "sqlite3", os.Getenv("WINTER"),
   )
   if e != nil {
      log.Fatal(e)
   }
   var mb string
   e = db.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist,
   ).Scan(&mb)
   if e != nil {
      log.Fatal(e)
   } else if mb == "" {
      log.Fatal("mb_s missing")
   }
   // local albums
   locals, e := localAlbum(db, mb)
   if e != nil {
      log.Fatal(e)
   }
   // remote albums
   remote, e := remoteAlbum(mb)
   if e != nil {
      log.Fatal(e)
   }
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
