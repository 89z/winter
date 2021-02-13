package main

import (
   "database/sql"
   "encoding/json"
   "fmt"
   "io/ioutil"
   "log"
   "os"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

type song struct {
   S string
}

func query(winter string) (*sql.Rows, error) {
   // get all artists
   db, e := sql.Open("sqlite3", winter)
   if e != nil {
      return nil, e
   }
   return db.Query("select artist_s from artist_t")
}

// find where the artist is not in the database
func main() {
   rows, e := query(os.Getenv("WINTER"))
   if e != nil {
      log.Fatal(e)
   }
   var (
      row string
      table = map[string]bool{}
   )
   for rows.Next() {
      e = rows.Scan(&row)
      if e != nil {
         log.Fatal(e)
      }
      table[strings.ToUpper(row)] = true
   }
   // JSON
   data, e := ioutil.ReadFile(os.Getenv("UMBER"))
   var songs []song
   e = json.Unmarshal(data, &songs)
   if e != nil {
      log.Fatal(e)
   }
   for _, each := range songs {
      artists := strings.Split(each.S, " - ")[0]
      artist := strings.Split(artists, ", ")[0]
      if ! table[strings.ToUpper(artist)] {
         fmt.Println(artist)
         table[strings.ToUpper(artist)] = true
      }
   }
}
