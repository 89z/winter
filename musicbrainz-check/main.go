package main

import (
   "database/sql"
   "errors"
   "fmt"
   "log"
   "os"
   "sort"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

type winterLocal struct {
   color string
   date string
}

func selectMb(db *sql.DB, artist string) (string, error) {
   var mb string
   e := db.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist,
   ).Scan(&mb)
   if e != nil {
      return "", e
   } else if mb == "" {
      return "", errors.New("mb_s missing")
   }
   return mb, nil
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
   mb, e := selectMb(db, artist)
   if e != nil {
      log.Fatal(e)
   }
   // local albums
   locals, e := localAlbum(db, mb)
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
         local, b := locals[strings.ToUpper(release)]
         if b {
            remotes[n].date = local.date
            remotes[n].color = local.color
         }
      }
   }
   sort.Slice(remotes, func(n, n2 int) bool {
      return remotes[n].date < remotes[n2].date
   })
   for _, group := range remotes {
      fmt.Printf("%-10v | %10v | %v\n", group.date, group.color, group.title)
   }
}
