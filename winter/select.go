package main

import (
   "database/sql"
   "fmt"
   "strings"
)

func SelectArtist(open_o *sql.DB, artist_s string) error {
   query_o, e := open_o.Query(`
      SELECT
         album_n,
         album_s,
         artist_n,
         check_s,
         date_s,
         note_s,
         song_n,
         song_s,
         url_s
      FROM album_t
      NATURAL JOIN song_album_t
      NATURAL JOIN song_t
      NATURAL JOIN song_artist_t
      NATURAL JOIN artist_t
      WHERE artist_s LIKE ?
      ORDER BY date_s
      `, artist_s,
   )
   if e != nil {
      return e
   }
   row_a := []Row{}
   pop_b := false
   for query_o.Next() {
      row := Row{}
      e = query_o.Scan(
         &row.Album,
         &row.AlbumStr,
         &row.Artist,
         &row.Check,
         &row.Date,
         &row.Note,
         &row.Song,
         &row.SongStr,
         &row.URL,
      )
      if e != nil {
         return e
      }
      if Pop(row.URL) {
         pop_b = true
      }
      row_a = append(row_a, row)
   }
   album_prev_n := 0
   for _, row := range row_a {
      if row.Album != album_prev_n {
         if album_prev_n > 0 {
            fmt.Println()
         }
         // print album number
         fmt.Println("album_n |", row.Album)
         // print album title
         fmt.Println("album_s |", row.AlbumStr)
         // print album date
         if row.Date != "" {
            fmt.Println("date_s  |", row.Date)
         } else {
            fmt.Println("date_s  |", YELLOW)
         }
         // print URL
         if pop_b {
            if row.URL != "" {
               fmt.Println("url_s   |", row.URL)
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
         album_prev_n = row.Album
      }
      // print song number, title
      fmt.Printf("%7v | %-*.*v | ", row.Song, WIDTH - 2, WIDTH - 2, row.SongStr)
      // print song note
      if row.Note == "" && ! Pop(row.URL) {
         fmt.Println(YELLOW)
      } else {
         fmt.Println(row.Note)
      }
   }
   fmt.Println()
   // print artist number
   fmt.Println("artist_n |", row_a[0].Artist)
   // print artist check
   if row_a[0].Check != "" {
      fmt.Println("check_s  |", row_a[0].Check)
   } else {
      fmt.Println("check_s  |", YELLOW)
   }
   return nil
}
