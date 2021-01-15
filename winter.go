package winter

import (
   "database/sql"
   "fmt"
   "strings"
)

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

