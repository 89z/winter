package main

import (
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
   check(e)
   id_s := o.Query().Get("v")
   info, e := youtube.Info(id_s)
   check(e)
   view_n, e := youtube.Views(info)
   check(e)
   color, _ := youtube.Color(view_n)
   println(color)
}
