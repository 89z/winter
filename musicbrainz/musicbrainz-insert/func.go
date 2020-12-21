package main

import (
   "database/sql"
   "fmt"
   "time"
)

func Note(f float64) string {
   d := time.Duration(f) * time.Millisecond
   if d < 179_500 * time.Millisecond {
      return "short"
   }
   if d > 15 * time.Minute {
      return "long"
   }
   return ""
}
