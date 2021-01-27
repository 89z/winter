package main

import (
   "encoding/json"
   "io/ioutil"
   "log"
   "strings"
)

type song struct {
   localId string
   year int
   remoteId string
   title string
}

func (s *song) UnmarshalJSON(b []byte) error {
   a := []interface{}{&s.localId, &s.year, &s.remoteId, &s.title}
   return json.Unmarshal(b, &a)
}

func main() {
   data, e := ioutil.ReadFile(`D:\Git\umber\umber.json`)
   if e != nil {
      log.Fatal(e)
   }
   var songs []song
   e = json.Unmarshal(data, &songs)
   if e != nil {
      log.Fatal(e)
   }
   done := map[string]bool{}
   for n, song := range songs {
      artist := strings.ToUpper(
         strings.Split(song.title, " - ")[0],
      )
      if n >= 30 && ! done[artist] {
         println(artist)
      }
      done[artist] = true
   }
}
