package main

import (
   "database/sql"
   "fmt"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

var (
   album_n int
   album_s string
   date_s string
   url_s string
   song_n int
   song_s string
   note_s string
   artist_n int
   check_s string
   pop_b bool
   should_b bool
)

func ArtistSelect(open_o *sql.DB, artist_s string) error {
   query_s := `
   SELECT
      album_n, album_s, date_s, url_s,
      song_n, song_s, note_s,
      artist_n, check_s, pop_n
   FROM album_t
   NATURAL JOIN song_album_t
   NATURAL JOIN song_t
   NATURAL JOIN song_artist_t
   NATURAL JOIN artist_t
   WHERE artist_s = ?
   ORDER BY date_s
   `
   query_o, e := open_o.Query(query_s, artist_s)
   if e != nil {
      return e
   }
   date_prev_s := "9999-12-31"
   for query_o.Next() {
      e = query_o.Scan(
         &album_n,
         &album_s,
         &date_s,
         &url_s,
         &song_n,
         &song_s,
         &note_s,
         &artist_n,
         &check_s,
         &pop_b,
      )
      if e != nil {
         return e
      }
      if date_s != date_prev_s {
         if date_prev_s != "9999-12-31" {
            fmt.Println()
         }
         // print album date
         if date_s != "" {
            fmt.Print(date_s)
         } else {
            fmt.Print("\x1b[30;43mdate ", album_n, "\x1b[m")
         }
         // print album title
         fmt.Println(" |", album_s)
         // print URL
         if pop_b {
            if url_s != "" {
               fmt.Println(url_s)
               should_b = true
            } else {
               fmt.Print("\x1b[30;43murl ", album_n, "\x1b[m\n")
            }
         }
         // print rule
         fmt.Println(strings.Repeat("-", 52))
         date_prev_s = date_s
      }
      // print song title
      fmt.Printf("%40.40v | ", song_s)
      // print song note
      if note_s == "" && ! strings.HasPrefix(url_s, "youtube.com/watch?") {
         fmt.Print("\x1b[30;43mnote ", song_n, "\x1b[m\n")
      } else {
         fmt.Println(note_s)
      }
   }
   fmt.Println()
   // print artist check
   if check_s != "" {
      fmt.Println("check:", check_s)
   } else {
      fmt.Print("\x1b[30;43mcheck ", artist_n, "\x1b[m\n")
   }
   // print artist pop
   if ! pop_b {
      fmt.Println("pop: false")
   } else if should_b {
      fmt.Println("pop: true")
   } else {
      fmt.Print("\x1b[30;43mpop ", artist_n, "\x1b[m\n")
   }
   return nil
}
