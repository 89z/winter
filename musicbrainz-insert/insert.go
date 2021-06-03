package main

import (
   "fmt"
   "github.com/89z/musicbrainz"
   "database/sql"
   "path"
   "strings"
   "time"
)

func insert(album musicbrainz.Release, tx *sql.Tx) error {
   // ALBUM
   result, err := tx.Exec(`
   INSERT INTO album_t (album_s, date_s, url_s) VALUES (?, ?, '')
   `, album.Title, album.Date)
   if err != nil { return err }
   albumId, err := result.LastInsertId()
   if err != nil { return err }
   // CREATE ARTIST ARRAY
   var artists []int
   for _, credit := range album.ArtistCredit {
      var artist int
      err := tx.QueryRow(`
      SELECT artist_n FROM artist_t WHERE mb_s = ?
      `, credit.Artist.ID).Scan(&artist)
      if err != nil {
         return fmt.Errorf("%v %v", credit.Name, err)
      }
      artists = append(artists, artist)
   }
   // CREATE SONG ARRAY
   var tns []titleNote
   for _, media := range album.Media {
      for _, track := range media.Tracks {
         tns = append(tns, titleNote{
            track.Title, note(track.Length),
         })
      }
   }
   // ITERATE SONG ARRAY
   for _, tn := range tns {
      result, err := tx.Exec(`
      INSERT INTO song_t (song_s, note_s, album_n) VALUES (?, ?, ?)
      `, tn.title, tn.note, albumId)
      if err != nil { return err }
      song, err := result.LastInsertId()
      if err != nil { return err }
      // ITERATE ARTIST ARRAY
      for _, artist := range artists {
         _, err := tx.Exec(`
         INSERT INTO song_artist_t VALUES (?, ?)
         `, song, artist)
         if err != nil { return err }
      }
   }
   return nil
}

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

type titleNote struct {
   title string
   note string
}

func release(addr string) (musicbrainz.Release, error) {
   id := path.Base(addr)
   if strings.Contains(addr, "musicbrainz.org/release/") {
      return musicbrainz.NewRelease(id)
   }
   g, err := musicbrainz.NewReleaseGroup(id)
   if err != nil {
      return musicbrainz.Release{}, err
   }
   g.Sort()
   return g.Releases[0], nil
}
