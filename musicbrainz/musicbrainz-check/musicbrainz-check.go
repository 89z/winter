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

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist_s := os.Args[1]
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   if e != nil {
      log.Fatal(e)
   }
   // local albums
   local_m, e := LocalAlbum(db, artist_s)
   if e != nil {
      log.Fatal(e)
   }
   // remote albums
   var mb_s string
   e = db.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist_s,
   ).Scan(&mb_s)
   if e != nil {
      log.Fatal(e)
   }
   if mb_s == "" {
      log.Fatal("mb_s missing")
   }
   remote_a, e := RemoteAlbum(mb_s)
   if e != nil {
      log.Fatal(e)
   }
   for n, group_o := range remote_a {
      for _, release_s := range group_o.Release {
         local_o, b := local_m[strings.ToUpper(release_s)]
         if b {
            remote_a[n].Date = local_o.Date
            remote_a[n].Color = local_o.Color
         }
      }
   }
   sort.Slice(remote_a, func(n, n2 int) bool {
      return remote_a[n].Date < remote_a[n2].Date
   })
   for _, group_o := range remote_a {
      fmt.Printf(
         "%-10v | %40.40v | %v\n", group_o.Date, group_o.Title, group_o.Color,
      )
   }
}
