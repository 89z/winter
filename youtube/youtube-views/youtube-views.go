package main

import (
   "log"
   "net/url"
   "os"
   "winter/youtube"
)

func main() {
   if len(os.Args) != 2 {
      println("youtube-views <URL>")
      os.Exit(1)
   }
   url_s := os.Args[1]
   o, e := url.Parse(url_s)
   if e != nil {
      log.Fatal(e)
   }
   id_s := o.Query().Get("v")
   m, e := youtube.Info(id_s)
   if e != nil {
      log.Fatal(e)
   }
   view_s, e := youtube.Views(m)
   if e != nil {
      log.Fatal(e)
   }
   println(view_s)
}
