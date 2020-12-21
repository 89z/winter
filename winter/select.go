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
   pop_n int
)

const YELLOW = "\x1b[43m   \x1b[m"
const WIDTH = 48

func SelectArtist(open_o *sql.DB, artist_s string) error {
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
   album_prev_n := 0
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
         &pop_n,
      )
      if e != nil {
         return e
      }
      if album_n != album_prev_n {
         if album_prev_n > 0 {
            fmt.Println()
         }
         // print album number
         fmt.Println("album_n |", album_n)
         // print album title
         fmt.Println("album_s |", album_s)
         // print album date
         if date_s != "" {
            fmt.Println("date_s  |", date_s)
         } else {
            fmt.Println("date_s  |", YELLOW)
         }
         // print URL
         if pop_n > 0 {
            if url_s != "" {
               fmt.Println("url_s   |", url_s)
            } else {
               fmt.Println("url_s   |", YELLOW)
            }
         }
         // print rule
         fmt.Print("--------+", strings.Repeat("-", WIDTH), "+-------\n")
         fmt.Print(
            "song_n  | song_s", strings.Repeat(" ", WIDTH - 7), "| note_s\n",
         )
         fmt.Print("--------+", strings.Repeat("-", WIDTH), "+-------\n")
         album_prev_n = album_n
      }
      // print song number, title
      fmt.Printf("%7v | %-*.*v | ", song_n, WIDTH - 2, WIDTH - 2, song_s)
      // print song note
      if note_s == "" && ! strings.HasPrefix(url_s, "youtube.com/watch?") {
         fmt.Println(YELLOW)
      } else {
         fmt.Println(note_s)
      }
   }
   fmt.Println()
   // print artist number
   fmt.Println("artist_n |", artist_n)
   // print artist pop
   fmt.Println("pop_n    |", pop_n)
   // print artist check
   if check_s != "" {
      fmt.Println("check_s  |", check_s)
   } else {
      fmt.Println("check_s  |", YELLOW)
   }
   return nil
}
