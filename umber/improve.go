package main

import (
   "database/sql"
   "encoding/json"
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
      artist string
      artists = map[string]bool{}
   )
   for query.Next() {
      e = query.Scan(&artist)
      x.Check(e)
      artists[strings.ToUpper(artist)] = true
   }
   // check JSON
   data, e := ioutil.ReadFile(os.Getenv("UMBER"))
   var songs []song
   e = json.Unmarshal(data, &songs)
   x.Check(e)
   for _, each := range songs {
      artist := strings.Split(each.S, " - ")[0]
      if ! artists[strings.ToUpper(artist)] {
         println(artist)
      }
   }
}
