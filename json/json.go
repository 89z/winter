package json
import "encoding/json"
type Map map[string]interface{}
type Slice []interface{}

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

func (a Slice) M(n int) Map {
   return a[n].(map[string]interface{})
}

func Decode(s string) (Map, error) {
   y := []byte(s)
   m := Map{}
   return m, json.Unmarshal(y, &m)
}

func Encode(a Slice) (string, error) {
   y, e := json.Marshal(a)
   if e != nil {
      return "", e
   }
   return string(y), nil
}
