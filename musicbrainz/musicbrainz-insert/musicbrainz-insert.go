package main

import (
   "github.com/89z/x"
   "log"
   "time"
)

func check(e error) {
   if e != nil {
      log.Fatal(e)
   }
}

func note(m x.Map) string {
   if m["length"] == nil {
      return "?:??"
   }
   n := time.Duration(m.N("length")) * time.Millisecond
   if n < 179_500 * time.Millisecond {
      return "short"
   }
   if n > 15 * time.Minute {
      return "long"
   }
   return ""
}

type song struct {
   title string
   note string
}
