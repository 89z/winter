package main

import (
   "database/sql"
   "fmt"
   "winter/snow"
)

func SelectArtist(open_o *sql.DB, artist_s string) error {
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
   // print artist number
   fmt.Println("artist_n |", artist_n)
   // print artist check
   if check_s != "" {
      fmt.Println("check_s  |", check_s)
   } else {
      fmt.Println("check_s  |", YELLOW)
   }
   // print musicbrainz id
   if mb_s != "" {
      fmt.Println("mb_s     |", mb_s)
   } else {
      fmt.Println("mb_s     |", YELLOW)
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
      if snow.Pop(r.UrlStr) {
         pop_b = true
      }
      row_a = append(row_a, r)
   }
   album_prev_n := 0
   for _, r := range row_a {
      if r.AlbumInt != album_prev_n {
         fmt.Println()
         // print album number
         fmt.Println("album_n |", r.AlbumInt)
         // print album title
         fmt.Println("album_s |", r.AlbumStr)
         // print album date
         if r.DateStr != "" {
            fmt.Println("date_s  |", r.DateStr)
         } else {
            fmt.Println("date_s  |", YELLOW)
         }
         // print URL
         if pop_b {
            if r.UrlStr != "" {
               fmt.Println("url_s   |", r.UrlStr)
            } else {
               fmt.Println("url_s   |", YELLOW)
            }
         }
         // print rule
         fmt.Print("--------+", DASH[:WIDTH], "+-------\n")
         fmt.Print("song_n  | song_s", SPACE[:WIDTH - 7], "| note_s\n")
         fmt.Print("--------+", DASH[:WIDTH], "+-------\n")
         album_prev_n = r.AlbumInt
      }
      // print song number, title
      fmt.Printf("%7v | %-*.*v | ", r.SongInt, WIDTH - 2, WIDTH - 2, r.SongStr)
      // print song note
      if r.NoteStr == "" && ! snow.Pop(r.UrlStr) {
         fmt.Println(YELLOW)
      } else {
         fmt.Println(r.NoteStr)
      }
   }
   return nil
}
