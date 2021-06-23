package main

import (
   "bytes"
   "database/sql"
   "fmt"
   "github.com/walles/moar/m"
   "strings"
)

const (
   reset = "\x1b[m"
   yellow = "\x1b[30;43m"
)

func selectOne(tx *sql.Tx, like string) error {
   // ARTIST
   var (
      artist, check, mb string
      artistId int
   )
   err := tx.QueryRow(`
   SELECT * FROM artist_t WHERE artist_s LIKE ?
   `, like).Scan(&artistId, &artist, &check, &mb)
   if err != nil {
      return err
   }
   // ALBUMS
   rows, err := tx.Query(`
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
   `, like)
   if err != nil {
      return err
   }
   defer rows.Close()
   var (
      records []record
      songs = make(map[string]int)
   )
   for rows.Next() {
      var r record
      err := rows.Scan(
         &r.albumInt,
         &r.albumStr,
         &r.dateStr,
         &r.noteStr,
         &r.songInt,
         &r.songStr,
         &r.urlStr,
      )
      if err != nil {
         return err
      }
      records = append(records, r)
      upper := strings.ToUpper(r.songStr)
      if songs[upper] == 0 {
         songs[upper] = 1
      } else {
         songs[upper]++
      }
   }
   b := new(bytes.Buffer)
   // print artist number
   fmt.Fprintln(b, "artist_n |", artistId)
   // print artist name
   fmt.Fprintln(b, "artist_s |", artist)
   // print artist check
   if check != "" {
      fmt.Fprintln(b, "check_s  |", check)
   } else {
      fmt.Fprintln(b, "check_s  |", yellow, " ", reset)
   }
   // print musicbrainz id
   if mb != "" {
      fmt.Fprintln(b, "mb_s     |", mb)
   } else {
      fmt.Fprintln(b, "mb_s     |", yellow, " ", reset)
   }
   var prev int
   for _, r := range records {
      if r.albumInt != prev {
         fmt.Fprintln(b)
         // print album number
         fmt.Fprintln(b, "album_n |", r.albumInt)
         // print album title
         fmt.Fprintln(b, "album_s |", r.albumStr)
         // print album date
         if r.dateStr != "" {
            fmt.Fprintln(b, "date_s  |", r.dateStr)
         } else {
            fmt.Fprintln(b, "date_s  |", yellow, " ", reset)
         }
         // print URL
         if r.urlStr != "" {
            fmt.Fprintln(b, "url_s   |", r.urlStr)
         } else {
            fmt.Fprintln(b, "url_s   |", yellow, " ", reset)
         }
         // print rule
         dash := strings.Repeat("-", 50)
         fmt.Fprint(b, "--------+-----------+", dash, "\n")
         fmt.Fprintln(b, "song_n  | note_s    | song_s")
         fmt.Fprint(b, "--------+-----------+", dash, "\n")
         prev = r.albumInt
      }
      // print song number
      fmt.Fprintf(b, "%7v | ", r.songInt)
      // print song note
      fmt.Fprint(b, note(r, songs), " | ")
      // print song title
      fmt.Fprintln(b, r.songStr)
   }
   read := m.NewReaderFromStream("winter", b)
   p := m.NewPager(read)
   p.DeInit, p.ShowLineNumbers = false, false
   return p.Page()
}

func note(r record, songs map[string]int) string {
   if r.noteStr != "" || strings.HasPrefix(r.urlStr, "youtube.com/watch?") {
      return fmt.Sprintf("%-9v", r.noteStr)
   }
   if songs[strings.ToUpper(r.songStr)] > 1 {
      return yellow + "duplicate" + reset
   }
   return yellow + "   " + reset + "      "
}

type record struct {
   albumInt int
   albumStr string
   dateStr string
   noteStr string
   songInt int
   songStr string
   urlStr string
}
