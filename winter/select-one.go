package main

import (
   "database/sql"
   "fmt"
   "os"
   "os/exec"
   "strings"
)

func SelectOne(open_o *sql.DB, artist_s string) error {
   // ARTIST
   var (
      artist_n int
      check_s string
      mb_s string
   )
   e := open_o.QueryRow(
      "select artist_n, check_s, mb_s from artist_t where artist_s LIKE ?",
      artist_s,
   ).Scan(&artist_n, &check_s, &mb_s)
   if e != nil {
      return e
   }
   // ALBUMS
   query_o, e := open_o.Query(`
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
   row_a := []Row{}
   song_m := map[string]int{}
   for query_o.Next() {
      r := Row{}
      e = query_o.Scan(
         &r.AlbumInt,
         &r.AlbumStr,
         &r.DateStr,
         &r.NoteStr,
         &r.SongInt,
         &r.SongStr,
         &r.UrlStr,
      )
      if e != nil {
         return e
      }
      row_a = append(row_a, r)
      upper := strings.ToUpper(r.SongStr)
      if song_m[upper] == 0 {
         song_m[upper] = 1
      } else {
         song_m[upper]++
      }
   }
   album_prev_n := 0
   less := exec.Command("less")
   pipe, e := less.StdinPipe()
   if e != nil {
      return e
   }
   less.Stdout = os.Stdout
   less.Start()
   /////////////////////////////////////////////////////////////////////////////
   // print artist number
   fmt.Fprintln(pipe, "artist_n |", artist_n)
   // print artist check
   if check_s != "" {
      fmt.Fprintln(pipe, "check_s  |", check_s)
   } else {
      fmt.Fprintln(pipe, "check_s  |", YELLOW)
   }
   // print musicbrainz id
   if mb_s != "" {
      fmt.Fprintln(pipe, "mb_s     |", mb_s)
   } else {
      fmt.Fprintln(pipe, "mb_s     |", YELLOW)
   }
   for _, r := range row_a {
      if r.AlbumInt != album_prev_n {
         fmt.Fprintln(pipe)
         // print album number
         fmt.Fprintln(pipe, "album_n |", r.AlbumInt)
         // print album title
         fmt.Fprintln(pipe, "album_s |", r.AlbumStr)
         // print album date
         if r.DateStr != "" {
            fmt.Fprintln(pipe, "date_s  |", r.DateStr)
         } else {
            fmt.Fprintln(pipe, "date_s  |", YELLOW)
         }
         // print URL
         if r.UrlStr != "" {
            fmt.Fprintln(pipe, "url_s   |", r.UrlStr)
         } else {
            fmt.Fprintln(pipe, "url_s   |", YELLOW)
         }
         // print rule
         fmt.Fprint(pipe, "--------+-----------+", DASH[:WIDTH], "\n")
         fmt.Fprintln(pipe, "song_n  | note_s    | song_s")
         fmt.Fprint(pipe, "--------+-----------+", DASH[:WIDTH], "\n")
         album_prev_n = r.AlbumInt
      }
      // print song number
      fmt.Fprintf(pipe, "%7v | ", r.SongInt)
      // print song note
      fmt_s, note_s := Note(r, song_m)
      fmt.Fprintf(pipe, fmt_s + " | ", note_s)
      // print song title
      fmt.Fprintln(pipe, r.SongStr)
   }
   /////////////////////////////////////////////////////////////////////////////
   pipe.Close()
   less.Wait()
   return nil
}