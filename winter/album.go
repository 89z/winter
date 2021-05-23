package main
import "database/sql"

func copyAlbum(tx *sql.Tx, source, dest string) error {
   var url string
   // COPY URL
   err := tx.QueryRow(`
   SELECT url_s FROM album_t WHERE album_n = ?
   `, source).Scan(&url)
   if err != nil { return err }
   // PASTE URL
   if _, err := tx.Exec(`
   UPDATE album_t SET url_s = ? WHERE album_n = ?
   `, url, dest); err != nil { return err }
   // COPY NOTES
   rows, err := tx.Query(`
   SELECT song_s, note_s FROM song_t WHERE album_n = ?
   `, source)
   if err != nil { return err }
   defer rows.Close()
   songs := make(map[string]string)
   for rows.Next() {
      var note, song string
      err := rows.Scan(&song, &note)
      if err != nil { return err }
      songs[song] = note
   }
   // PASTE NOTES
   for song, note := range songs {
      _, err := tx.Exec(`
      UPDATE song_t SET note_s = ?
      WHERE album_n = ? AND song_s = ? COLLATE NOCASE
      `, note, dest, song)
      if err != nil { return err }
   }
   return nil
}

func deleteAlbum(tx *sql.Tx, album string) error {
   rows, err := tx.Query("SELECT song_n FROM song_t WHERE album_n = ?", album)
   if err != nil { return err }
   defer rows.Close()
   var songs []int
   for rows.Next() {
      var song int
      err := rows.Scan(&song)
      if err != nil { return err }
      songs = append(songs, song)
   }
   for _, song := range songs {
      _, err := tx.Exec("DELETE FROM song_t WHERE song_n = ?", song)
      if err != nil { return err }
      if _, err := tx.Exec("DELETE FROM song_artist_t WHERE song_n = ?", song)
      err != nil { return err }
   }
   {
      _, err := tx.Exec("DELETE FROM album_t WHERE album_n = ?", album)
      return err
   }
}
