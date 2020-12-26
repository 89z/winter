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
   db_s := os.Getenv("WINTER")
   open_o, e := sql.Open("sqlite3", db_s)
   if e != nil {
      log.Fatal(e)
   }
   // local albums
   local_m, e := LocalAlbum(open_o, artist_s)
   if e != nil {
      log.Fatal(e)
   }
   // remote albums
   var mb_s string
   e = open_o.QueryRow(
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
   for n, album_o := range remote_a {
      local_o, b := local_m[strings.ToUpper(album_o.Title)]
      if b {
         remote_a[n].Date = local_o.Date
      }
   }
   sort.Slice(remote_a, func(n, n2 int) bool {
      return remote_a[n].Date < remote_a[n2].Date
   })
   remote_m := map[string]bool{}
   for _, album_o := range remote_a {
      if remote_m[album_o.Group] {
         continue
      }
      remote_m[album_o.Group] = true
      color_s := local_m[strings.ToUpper(album_o.Title)].Color
      fmt.Printf("%-10v | %40.40v | %v\n", album_o.Date, album_o.Title, color_s)
   }
}
