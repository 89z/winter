package dbstorm

import (
   "fmt"
   "github.com/asdine/storm/v3"
)

func Save(tx storm.Node, row interface{}) error {
   e := tx.Save(row)
   if e != nil {
      return e
   }
   fmt.Printf("%+v\n", row)
   return nil
}

func Update(tx storm.Node, row interface{}) error {
   e := tx.Update(row)
   if e != nil {
      return e
   }
   fmt.Printf("%+v\n", row)
   return nil
}
