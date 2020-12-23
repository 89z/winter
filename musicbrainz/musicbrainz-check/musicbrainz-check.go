package main

import (
   "database/sql"
   "fmt"
   "log"
   "os"
   "sort"
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
   sort.Slice(remote_a, func(n, n2 int) bool {
      return remote_a[n][0] < remote_a[n2][0]
   })
   remote_m := map[string]bool{}
   for _, album_a := range remote_a {
      album_s := album_a[1]
      if remote_m[album_s] {
         continue
      }
      remote_m[album_s] = true
      date_s := album_a[0]
      color_s := local_m[album_s]
      fmt.Printf("%-10v | %40.40v | %v\n", date_s, album_s, color_s)
   }
}
