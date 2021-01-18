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

func check(e error) {
   if e != nil {
      log.Fatal(e)
   }
}

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist_s := os.Args[1]
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   check(e)
   // local albums
   local_m, e := localAlbum(db, artist_s)
   check(e)
   // remote albums
   var mb_s string
   e = db.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist_s,
   ).Scan(&mb_s)
   check(e)
   if mb_s == "" {
      log.Fatal("mb_s missing")
   }
   remote_a, e := remoteAlbum(mb_s)
   check(e)
   for n, group := range remote_a {
      for release_s := range group.release {
         local_o, b := local_m[strings.ToUpper(release_s)]
         if b {
            remote_a[n].date = local_o.date
            remote_a[n].color = local_o.color
         }
      }
   }
   sort.Slice(remote_a, func(n, n2 int) bool {
      return remote_a[n].date < remote_a[n2].date
   })
   for _, group := range remote_a {
      fmt.Printf(
         "%-10v | %40.40v | %v\n", group.date, group.title, group.color,
      )
   }
}
