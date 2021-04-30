package main

import (
   "fmt"
   "github.com/89z/rosso/musicbrainz"
   "time"
   "winter"
)

func note(length int) string {
   if length == 0 {
      return "?:??"
   }
   dur := time.Duration(length) * time.Millisecond
   if dur < 179_500 * time.Millisecond {
      return "short"
   }
   if dur > 15 * time.Minute {
      return "long"
   }
   return ""
}

type titleNote struct { title, note string }

func insert(album musicbrainz.Release, tx winter.Tx) error {
   // ALBUM
   albumId, e := tx.Insert(
      "album_t (album_s, date_s, url_s) values (?, ?, '')",
      album.Title,
      album.Date,
   )
   if e != nil { return e }
   // CREATE ARTIST ARRAY
   var artists []int
   for _, each := range album.ArtistCredit {
      var artist int
      e = tx.QueryRow(
         "select artist_n from artist_t where mb_s = ?", each.Artist.Id,
      ).Scan(&artist)
      if e != nil {
         return fmt.Errorf("%v %v", each.Name, e)
      }
      artists = append(artists, artist)
   }
   // CREATE SONG ARRAY
   var songs []titleNote
   for _, media := range album.Media {
      for _, track := range media.Tracks {
         songs = append(songs, titleNote{
            track.Title, note(track.Length),
         })
      }
   }
   // ITERATE SONG ARRAY
   for _, each := range songs {
      song, e := tx.Insert(
         "song_t (song_s, note_s, album_n) values (?, ?, ?)",
         each.title,
         each.note,
         albumId,
      )
      if e != nil { return e }
      // ITERATE ARTIST ARRAY
      for _, artist := range artists {
         _, e = tx.Insert("song_artist_t values (?, ?)", song, artist)
         if e != nil { return e }
      }
   }
   return nil
}
