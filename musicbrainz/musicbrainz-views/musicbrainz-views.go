package main

import (
   "io/ioutil"
   "log"
   "net/http"
   "net/url"
   "regexp"
)

func check(e error) {
   if e != nil {
      log.Fatal(e)
   }
}

func findSubmatch(re string, input []byte) string {
   a := regexp.MustCompile(re).FindSubmatch(input)
   if len(a) < 2 {
      return ""
   }
   return string(a[1])
}

func getContents(s string) ([]byte, error) {
   o, e := http.Get(s)
   if e != nil {
      return []byte{}, e
   }
   return ioutil.ReadAll(o.Body)
}

func youtubeResult(query string) (string, error) {
   value := url.Values{}
   value.Set("search_query", query)
   result := "https://www.youtube.com/results?" + value.Encode()
   println(result)
   get, e := getContents(result)
   if e != nil {
      return "", e
   }
   return findSubmatch("/vi/([^/]*)/", get), nil
}
