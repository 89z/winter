package main

import (
   "io/ioutil"
   "net/http"
   "net/url"
   "regexp"
)

func GetContents(s string) ([]byte, error) {
   o, e := http.Get(s)
   if e != nil {
      return []byte{}, e
   }
   return ioutil.ReadAll(o.Body)
}

func FindSubmatch(pat string, sub []byte) []byte {
   a := regexp.MustCompile(pat).FindSubmatch(sub)
   if len(a) < 2 {
      return []byte{}
   }
   return a[1]
}

func YoutubeResult(query_s string) (string, error) {
   m := url.Values{}
   m.Set("search_query", query_s)
   res_s := "https://www.youtube.com/results?" + m.Encode()
   println(res_s)
   get_y, e := GetContents(res_s)
   if e != nil {
      return "", e
   }
   find_y := FindSubmatch("/vi/([^/]*)/", get_y)
   return string(find_y), nil
}
