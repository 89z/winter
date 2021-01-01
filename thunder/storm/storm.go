package dbstorm

import (
   "github.com/asdine/storm/v3"
   "log"
   "os"
)

func main() {
   if len(os.Args) != 2 {
      println(`Save artist:
   thunder artist 'Kate Bush'

Select one artist:
   thunder 'Kate Bush'`)
      os.Exit(1)
   }
   db_s := os.Getenv("THUNDER")
   db_o, e := storm.Open(db_s)
   if e != nil {
      log.Fatal(e)
   }
   key_s := os.Args[1]
   switch key_s {
   case "artist":
      artist_s := os.Args[2]
      artist_o := Artist{Name: artist_s}
      e = Save(db_o, &artist_o)
   default:
      e = SelectOne(db_o, key_s)
   }
   if e != nil {
      log.Fatal(e)
   }
}
