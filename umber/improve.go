package main

import (
   "encoding/json"
   "fmt"
   "log"
   "os"
   "strings"
   "winter"
)

type song struct {
   S string
}

func query(file string) (*winter.Rows, error) {
   // get all artists
   tx, e := winter.NewTx(file)
   if e != nil {
      return nil, e
   }
   return tx.Query("select artist_s from artist_t")
}

// find where the artist is not in the database
func main() {
   rows, e := query(
      os.Getenv("WINTER"),
   )
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
   data, e := os.ReadFile(os.Getenv("UMBER"))
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
