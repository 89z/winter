package winter

import (
   "encoding/json"
   "net/http"
   "strings"
)

func JsonGetHttp(s string) (Map, error) {
   println(s)
   o, e := http.Get(s)
   if e != nil {
      return nil, e
   }
   m := Map{}
   return m, json.NewDecoder(o.Body).Decode(&m)
}

func Pop(s string) bool {
   return strings.HasPrefix(s, "youtube.com/watch?")
}
