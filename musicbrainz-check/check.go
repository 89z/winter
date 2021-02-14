package main

import (
   "errors"
   "fmt"
   "log"
   "os"
   "sort"
   "strings"
   "strings"
   "winter"
   "winter"
)

type winterLocal struct {
   color string
   date string
}

func selectMb(tx winter.Tx, artist string) (string, error) {
   var mb string
   e := tx.QueryRow(
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

func color(url string, unrated, good int) string {
   const (
      block = "\u2587\u2587\u2587\u2587\u2587"
      greenFive = "\x1b[92m" + block + "\x1b[90m" + block + "\x1b[m"
      greenTen = "\x1b[92m" + block + block + "\x1b[m"
      redFive = "\x1b[91m" + block + "\x1b[90m" + block + "\x1b[m"
      redTen = "\x1b[91m" + block + block + "\x1b[m"
   )
   if strings.HasPrefix(url, "youtube.com/watch?") {
      return greenTen
   }
   if unrated == 0 && good == 0 {
      return redTen
   }
   if unrated == 0 {
      return greenTen
   }
   if good == 0 {
      return redFive
   }
   return greenFive
}

func localAlbum(tx winter.Tx, mb string) (map[string]winterLocal, error) {
   query, e := tx.Query(`
   select
      album_s,
      date_s,
      url_s,
      count(1) filter (where note_s = '') as unrated,
      count(1) filter (where note_s = 'good') as good
   from album_t
   natural join song_t
   natural join song_artist_t
   natural join artist_t
   where mb_s = ?
   group by album_n
   `, mb)
   if e != nil {
      return nil, e
   }
   var (
      locals = map[string]winterLocal{}
      q queryRow
   )
   for query.Next() {
      e = query.Scan(&q.album, &q.date, &q.url, &q.unrated, &q.good)
      if e != nil {
         return nil, e
      }
      locals[strings.ToUpper(q.album)] = winterLocal{
         color(q.url, q.unrated, q.good), q.date,
      }
   }
   return locals, nil
}

type queryRow struct {
   album string
   date string
   url string
   unrated int
   good int
}
