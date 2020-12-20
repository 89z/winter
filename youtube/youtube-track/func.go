package main

import (
   "net/http"
   "regexp"
)

func FindSubmatch(pat, sub string) string {
   a := regexp.MustCompile(pat).FindStringSubmatch(sub)
   if len(a) < 2 {
      return ""
   }
   return a[1]
}

func GetImage(id_s string) string {
   url_s := "https://i.ytimg.com/vi/"
   if HttpHead(url_s + id_s + "/sddefault.jpg") {
      return ""
   }
   if HttpHead(url_s + id_s + "/sd1.jpg") {
      return "/sd1"
   }
   return "/hqdefault"
}

func HttpHead(s string) bool {
   println(s)
   o, e := http.Head(s)
   return e == nil && o.StatusCode == 200
}
