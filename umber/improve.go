package main

import (
   "database/sql"
   "encoding/json"
   "fmt"
   "os"
   "strings"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   db, err := sql.Open("sqlite3", os.Getenv("WINTER"))
   if err != nil {
      panic(err)
   }
   defer db.Close()
   rows, err := db.Query("SELECT artist_s FROM artist_t")
   if err != nil {
      panic(err)
   }
   defer rows.Close()
   table := make(map[string]bool)
   for rows.Next() {
      var row string
      err := rows.Scan(&row)
      if err != nil {
         panic(err)
      }
      table[strings.ToUpper(row)] = true
   }
   data, err := os.ReadFile(os.Getenv("UMBER"))
   if err != nil {
      panic(err)
   }
   var songs []struct { S string }
   err = json.Unmarshal(data, &songs)
   if err != nil {
      panic(err)
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
