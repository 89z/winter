package main

import (
   "fmt"
   "github.com/89z/rosso/musicbrainz"
   "os"
   "winter"
)

func main() {
   if len(os.Args) != 2 {
      fmt.Println(`musicbrainz-insert <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      return
   }
   tx, err := winter.NewTx(os.Getenv("WINTER"))
   if err != nil {
      panic(err)
   }
   album, err := musicbrainz.NewRelease(os.Args[1])
   if err != nil {
      panic(err)
   }
   err = insert(album, tx)
   if err != nil {
      panic(err)
   }
   err = tx.Commit()
   if err != nil {
      panic(err)
   }
}
