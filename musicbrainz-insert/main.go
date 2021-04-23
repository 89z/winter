package main

import (
   "fmt"
   "github.com/89z/x/musicbrainz"
   "os"
   "winter"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   tx, e := winter.NewTx(os.Getenv("WINTER"))
   if e != nil {
      panic(e)
   }
   album, e := musicbrainz.NewRelease(os.Args[1])
   if e != nil {
      panic(e)
   }
   e = insert(album, tx)
   if e != nil {
      panic(e)
   }
   e = tx.Commit()
   if e != nil {
      panic(e)
   }
}
