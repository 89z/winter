package main

import (
   "time"
   "winter"
)

func Note(m winter.Map) string {
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

type Song struct {
   Title string
   Note string
}
