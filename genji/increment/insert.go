package main

import (
   "github.com/genjidb/genji"
   "github.com/genjidb/genji/document"
   "log"
)

type Artist struct {
   ID int
   Name string
}

func check(e error) {
   if e != nil {
      log.Fatal(e)
   }
}

func main() {
   db, e := genji.Open("thunder.db")
   check(e)
   defer db.Close()
   e = db.Exec("CREATE TABLE artist_t")
   check(e)
   tx, e := db.DB.Begin(true)
   check(e)
   tb, e := tx.GetTable("artist_t")
   check(e)
   // 1
   d1, e := document.NewFromStruct(Artist{Name: "Sting"})
   check(e)
   _, e = tb.Insert(d1)
   check(e)
   // 2
   d2, e := document.NewFromStruct(Artist{Name: "Madonna"})
   check(e)
   _, e = tb.Insert(d2)
   check(e)
   tx.Commit()
   // query
   res, e := db.Query("SELECT id, name FROM artist_t")
   check(e)
   defer res.Close()
   e = res.Iterate(func (d document.Document) error {
      var (
         id int
         name string
      )
      e = document.Scan(d, &id, &name)
      if e != nil {
         return e
      }
      println(id, name)
      return nil
   })
   check(e)
}
