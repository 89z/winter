package main
import "winter"

func copyAlbum(tx winter.Tx, source , dest string) error {
   var note, song, url string
   // COPY URL
   e := tx.QueryRow(
      "select url_s from album_t where album_n = ?", source,
   ).Scan(&url)
   if e != nil {
      return e
   }
   // PASTE URL
   e = tx.Update("album_t set url_s = ? where album_n = ?", url, dest)
   if e != nil {
      return e
   }
   // COPY NOTES
   query, e := tx.Query(
      "select song_s, note_s from song_t where album_n = ?", source,
   )
   if e != nil {
      return e
   }
   songs := make(map[string]string)
   for query.Next() {
      e = query.Scan(&song, &note)
      if e != nil {
         return e
      }
      songs[song] = note
   }
   // PASTE NOTES
   for song, note := range songs {
      e = tx.Update(`
      song_t set note_s = ?
      where album_n = ? and song_s = ? COLLATE NOCASE
      `, note, dest, song)
      if e != nil {
         return e
      }
   }
   return nil
}

func deleteAlbum(tx winter.Tx, album string) error {
   query, e := tx.Query("select song_n from song_t where album_n = ?", album)
   if e != nil {
      return e
   }
   var (
      song int
      songs []int
   )
   for query.Next() {
      e = query.Scan(&song)
      if e != nil {
         return e
      }
      songs = append(songs, song)
   }
   for _, song := range songs {
      e = tx.Delete("song_t where song_n = ?", song)
      if e != nil {
         return e
      }
      e = tx.Delete("song_artist_t where song_n = ?", song)
      if e != nil {
         return e
      }
   }
   e = tx.Delete("album_t where album_n = ?", album)
   if e != nil {
      return e
   }
   return nil
}
