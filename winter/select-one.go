package main

import (
   "bytes"
   "fmt"
   "github.com/89z/page"
   "io"
   "strings"
   "winter"
)

func selectOne(tx winter.Tx, like string) error {
   // ARTIST
   var (
      artist, check, mb string
      artistId int
   )
   e := tx.QueryRow(
      "select * from artist_t where artist_s LIKE ?", like,
   ).Scan(&artistId, &artist, &check, &mb)
   if e != nil { return e }
   // ALBUMS
   query, e := tx.Query(`
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
   if e != nil { return e }
   var (
      rows []row
      songs = make(map[string]int)
   )
   for query.Next() {
      var r row
      e = query.Scan(
         &r.albumInt,
         &r.albumStr,
         &r.dateStr,
         &r.noteStr,
         &r.songInt,
         &r.songStr,
         &r.urlStr,
      )
      if e != nil { return e }
      rows = append(rows, r)
      upper := strings.ToUpper(r.songStr)
      if songs[upper] == 0 {
         songs[upper] = 1
      } else {
         songs[upper]++
      }
   }
   var (
      b = new(bytes.Buffer)
      prev = 0
   )
   // print artist number
   fmt.Fprintln(b, "artist_n |", artistId)
   // print artist name
   fmt.Fprintln(b, "artist_s |", artist)
   // print artist check
   if check != "" {
      fmt.Fprintln(b, "check_s  |", check)
   } else {
      fmt.Fprintln(b, "check_s  |", yellow)
   }
   // print musicbrainz id
   if mb != "" {
      fmt.Fprintln(b, "mb_s     |", mb)
   } else {
      fmt.Fprintln(b, "mb_s     |", yellow)
   }
   for _, r := range rows {
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
            fmt.Fprintln(b, "date_s  |", yellow)
         }
         // print URL
         if r.urlStr != "" {
            fmt.Fprintln(b, "url_s   |", r.urlStr)
         } else {
            fmt.Fprintln(b, "url_s   |", yellow)
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
      format, songNote := note(r, songs)
      fmt.Fprintf(b, format + " | ", songNote)
      // print song title
      fmt.Fprintln(b, r.songStr)
   }
   doc, e := page.NewDocument()
   if e != nil { return e }
   doc.ReadAll(io.NopCloser(b))
   root, e := page.NewOviewer(doc)
   if e != nil { return e }
   root.Run()
   root.WriteOriginal()
   return nil
}

func note(r row, songs map[string]int) (string, string) {
   switch {
   case r.noteStr != "", strings.HasPrefix(r.urlStr, "youtube.com/watch?"):
      return "%-9v", r.noteStr
   case songs[strings.ToUpper(r.songStr)] > 1:
      return "\x1b[30;43m%v\x1b[m", "duplicate"
   default:
      return yellow + "%6v", ""
   }
}

const yellow = "\x1b[43m   \x1b[m"

type row struct {
   albumInt, songInt int
   albumStr, dateStr, noteStr, songStr,urlStr string
}
