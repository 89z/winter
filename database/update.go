package main

import (
   "database/sql"
   "fmt"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

func CheckUpdate(open_o *sql.DB, artist_s, check_s string) error {
   query_s := `
   UPDATE artist_t SET check_s = ?
   WHERE artist_n = ?
   `
   exec_o, e := open_o.Exec(query_s, check_s, artist_s)
   if e != nil {
      return fmt.Errorf("%v %v", exec_o, e)
   }
   return nil
}

func DateUpdate(open_o *sql.DB, album_s, date_s string) error {
   query_s := `
   UPDATE album_t SET date_s = ?
   WHERE album_n = ?
   `
   exec_o, e := open_o.Exec(query_s, date_s, album_s)
   if e != nil {
      return fmt.Errorf("%v %v", exec_o, e)
   }
   return nil
}

func NoteUpdate(open_o *sql.DB, song_s, note_s string) error {
   query_s := `
   UPDATE song_t SET note_s = ?
   WHERE song_n = ?
   `
   exec_o, e := open_o.Exec(query_s, note_s, song_s)
   if e != nil {
      return fmt.Errorf("%v %v", exec_o, e)
   }
   return nil
}

func PopUpdate(open_o *sql.DB, artist_s, pop_s string) error {
   query_s := `
   UPDATE artist_t SET pop_n = ?
   WHERE artist_n = ?
   `
   exec_o, e := open_o.Exec(query_s, pop_s, artist_s)
   if e != nil {
      return fmt.Errorf("%v %v", exec_o, e)
   }
   return nil
}

func UrlUpdate(open_o *sql.DB, album_s, url_s string) error {
   query_s := `
   UPDATE album_t SET url_s = ?
   WHERE album_n = ?
   `
   exec_o, e := open_o.Exec(query_s, url_s, album_s)
   if e != nil {
      return fmt.Errorf("%v %v", exec_o, e)
   }
   return nil
}
