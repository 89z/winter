package main

import (
   "database/sql"
   "encoding/json"
   "fmt"
   "github.com/89z/x"
   "io/ioutil"
   "os"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

type song struct {
   S string
}

// find rows where the artist is not in the database
func main() {
   // get all artists
   db, e := sql.Open("sqlite3", os.Getenv("WINTER"))
   x.Check(e)
   query, e := db.Query("select artist_s from artist_t")
   x.Check(e)
   var (
      row string
      table = make(map[string]bool)
   )
   for query.Next() {
      e = query.Scan(&row)
      x.Check(e)
      table[strings.ToUpper(row)] = true
   }
   // check JSON
   data, e := ioutil.ReadFile(os.Getenv("UMBER"))
   var songs []song
   e = json.Unmarshal(data, &songs)
   x.Check(e)
   for _, each := range songs {
      artists := strings.Split(each.S, " - ")[0]
      artist := strings.Split(artists, ", ")[0]
      if ! table[strings.ToUpper(artist)] {
         fmt.Println(artist)
         table[strings.ToUpper(artist)] = true
      }
   }
}
