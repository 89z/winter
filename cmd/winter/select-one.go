package main

import (
   "database/sql"
   "fmt"
   "strings"
   "winter"
)

const (
   dash = "-----------------------------------------------------------------"
   width = 50
   yellow = "\x1b[43m   \x1b[m"
)

func note(r row, song_m map[string]int) (string, string) {
   if r.noteStr != "" || winter.Pop(r.urlStr) {
      return "%-9v", r.noteStr
   }
   if song_m[strings.ToUpper(r.songStr)] > 1 {
      return "\x1b[30;43m%v\x1b[m", "duplicate"
   }
   return yellow + "%6v", ""
}

func selectOne(tx *sql.Tx, artist_s string) error {
   // ARTIST
   var (
      artist_n int
      check_s string
      mb_s string
   )
   e := tx.QueryRow(
      "select artist_n, check_s, mb_s from artist_t where artist_s LIKE ?",
      artist_s,
   ).Scan(&artist_n, &check_s, &mb_s)
   if e != nil {
      return e
   }
   // ALBUMS
   query_o, e := tx.Query(`
      SELECT
         album_n,
         album_s,
         date_s,
         note_s,
         song_n,
         song_s,
         url_s
      FROM album_t
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
   row_a := []row{}
   song_m := map[string]int{}
   for query_o.Next() {
      r := row{}
      e = query_o.Scan(
         &r.albumInt,
         &r.albumStr,
         &r.dateStr,
         &r.noteStr,
         &r.songInt,
         &r.songStr,
         &r.urlStr,
      )
      if e != nil {
         return e
      }
      row_a = append(row_a, r)
      upper := strings.ToUpper(r.songStr)
      if song_m[upper] == 0 {
         song_m[upper] = 1
      } else {
         song_m[upper]++
      }
   }
   album_prev_n := 0
   cmd, pipe, e := less()
   if e != nil {
      return e
   }
   defer cmd.Wait()
   defer pipe.Close()
   // print artist number
   fmt.Fprintln(pipe, "artist_n |", artist_n)
   // print artist check
   if check_s != "" {
      fmt.Fprintln(pipe, "check_s  |", check_s)
   } else {
      fmt.Fprintln(pipe, "check_s  |", yellow)
   }
   // print musicbrainz id
   if mb_s != "" {
      fmt.Fprintln(pipe, "mb_s     |", mb_s)
   } else {
      fmt.Fprintln(pipe, "mb_s     |", yellow)
   }
   for _, r := range row_a {
      if r.albumInt != album_prev_n {
         fmt.Fprintln(pipe)
         // print album number
         fmt.Fprintln(pipe, "album_n |", r.albumInt)
         // print album title
         fmt.Fprintln(pipe, "album_s |", r.albumStr)
         // print album date
         if r.dateStr != "" {
            fmt.Fprintln(pipe, "date_s  |", r.dateStr)
         } else {
            fmt.Fprintln(pipe, "date_s  |", yellow)
         }
         // print URL
         if r.urlStr != "" {
            fmt.Fprintln(pipe, "url_s   |", r.urlStr)
         } else {
            fmt.Fprintln(pipe, "url_s   |", yellow)
         }
         // print rule
         fmt.Fprint(pipe, "--------+-----------+", dash[:width], "\n")
         fmt.Fprintln(pipe, "song_n  | note_s    | song_s")
         fmt.Fprint(pipe, "--------+-----------+", dash[:width], "\n")
         album_prev_n = r.albumInt
      }
      // print song number
      fmt.Fprintf(pipe, "%7v | ", r.songInt)
      // print song note
      fmt_s, note_s := note(r, song_m)
      fmt.Fprintf(pipe, fmt_s + " | ", note_s)
      // print song title
      fmt.Fprintln(pipe, r.songStr)
   }
   return nil
}

type row struct {
   albumInt int
   albumStr string
   dateStr string
   noteStr string
   songInt int
   songStr string
   urlStr string
}
