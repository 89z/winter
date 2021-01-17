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

func JsonGetString(s string) (Map, error) {
   y := []byte(s)
   m := Map{}
   return m, json.Unmarshal(y, &m)
}

func Pop(s string) bool {
   return strings.HasPrefix(s, "youtube.com/watch?")
}

type Map map[string]interface{}

func (m Map) A(s string) Slice {
   return m[s].([]interface{})
}

func (m Map) M(s string) Map {
   return m[s].(map[string]interface{})
}

func (m Map) N(s string) float64 {
   return m[s].(float64)
}

func (m Map) S(s string) string {
   return m[s].(string)
}

type Slice []interface{}

func (a Slice) M(n int) Map {
   return a[n].(map[string]interface{})
}
