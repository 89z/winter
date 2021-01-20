package youtube

import (
   "fmt"
   "github.com/89z/sienna/json"
   "io/ioutil"
   "math"
   "net/http"
   "net/url"
   "strconv"
   "time"
)

func Color(n float64) (string, bool) {
   s := numberFormat(n)
   if n > 8_000_000 {
      return "\x1b[1;31m" + s + "\x1b[m", true
   }
   return "\x1b[1;32m" + s + "\x1b[m", false
}

func floatVal(s string) (float64, error) {
   return strconv.ParseFloat(s, 64)
}

func getContents(s string) (string, error) {
   println(s)
   o, e := http.Get(s)
   if e != nil {
      return "", e
   }
   y, e := ioutil.ReadAll(o.Body)
   if e != nil {
      return "", e
   }
   return string(y), nil
}

func Info(id string) (json.Map, error) {
   info := "https://www.youtube.com/get_video_info?video_id=" + id
   query, e := getContents(info)
   if e != nil {
      return nil, e
   }
   v, e := url.ParseQuery(query)
   if e != nil {
      return nil, e
   }
   resp := v.Get("player_response")
   m, e := json.Load(resp)
   if e != nil {
      return nil, e
   }
   return m.M("microformat").M("playerMicroformatRenderer"), nil
}

func numberFormat(n float64) string {
   n2 := int(math.Log10(n)) / 3
   n3 := n / math.Pow10(n2 * 3)
   return fmt.Sprintf("%.3f", n3) + []string{"", " k", " M", " B"}[n2]
}

func sinceHours(left string) (float64, error) {
   right := "1970-01-01T00:00:00Z"[len(left):]
   o, e := time.Parse(time.RFC3339, left + right)
   if e != nil {
      return 0, e
   }
   return time.Since(o).Hours(), nil
}

func Views(m json.Map) (float64, error) {
   view_s := m.S("viewCount")
   view_n, e := floatVal(view_s)
   if e != nil {
      return 0, e
   }
   date_s := m.S("publishDate")
   hour_n, e := sinceHours(date_s)
   if e != nil {
      return 0, e
   }
   return view_n / (hour_n / 24 / 365), nil
}
