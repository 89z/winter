package main

import (
   "fmt"
   "github.com/89z/x/musicbrainz"
   "log"
   "os"
   "time"
   "winter"
)

type titleNote struct {
   title, note string
}

func note(length int) string {
   if length == 0 {
      return "?:??"
   }
   dur := time.Duration(length) * time.Millisecond
   if dur < 179_500 * time.Millisecond {
      return "short"
   }
   if dur > 15 * time.Minute {
      return "long"
   }
   return ""
}
