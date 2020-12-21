package main
import "database/sql"

func InsertArtist(open_o *sql.DB, artist_s string) error {
   _, e := Insert(
      open_o,
      "artist_t (artist_s, check_s, pop_n) values (?, '', 1)",
      artist_s,
   )
   return e
}
