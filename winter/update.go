package main
import "database/sql"

func UpdateCheck(open_o *sql.DB, artist_s, check_s string) error {
   _, e := open_o.Exec(
      "update artist_t set check_s = ? where artist_n = ?", check_s, artist_s,
   )
   return e
}

func UpdateDate(open_o *sql.DB, album_s, date_s string) error {
   _, e := open_o.Exec(
      "update album_t set date_s = ? where album_n = ?", date_s, album_s,
   )
   return e
}

func UpdateNote(open_o *sql.DB, song_s, note_s string) error {
   _, e := open_o.Exec(
      "update song_t set note_s = ? where song_n = ?", note_s, song_s,
   )
   return e
}

func UpdatePop(open_o *sql.DB, artist_s, pop_s string) error {
   _, e := open_o.Exec(
      "update artist_t set pop_n = ? where artist_n = ?", pop_s, artist_s,
   )
   return e
}

func UpdateURL(open_o *sql.DB, album_s, url_s string) error {
   _, e := open_o.Exec(
      "update album_t set url_s = ? where album_n = ?", url_s, album_s,
   )
   return e
}
