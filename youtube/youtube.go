package youtube

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
   "math"
   "net/http"
   "net/url"
   "strconv"
   "time"
   "winter"
)

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

func numberFormat(n float64) string {
   n2 := int(math.Log10(n)) / 3
   n3 := n / math.Pow10(n2 * 3)
   return fmt.Sprintf("%.3f", n3) + []string{"", " k", " M", " B"}[n2]
}

func Color(n float64) (string, bool) {
   s := numberFormat(n)
   if n > 8_000_000 {
      return "\x1b[1;31m" + s + "\x1b[m", true
   }
   return "\x1b[1;32m" + s + "\x1b[m", false
}

func Info(id_s string) (winter.Map, error) {
   info_s := "https://www.youtube.com/get_video_info?video_id=" + id_s
   query_s, e := getContents(info_s)
   if e != nil {
      return nil, e
   }
   o, e := url.ParseQuery(query_s)
   if e != nil {
      return nil, e
   }
   resp_s := o.Get("player_response")
   json_m := winter.Map{}
   e = json.Unmarshal([]byte(resp_s), &json_m)
   if e != nil {
      return nil, e
   }
   return json_m.M("microformat").M("playerMicroformatRenderer"), nil
}

func sinceHours(left string) (float64, error) {
   right := "1970-01-01T00:00:00Z"[len(left):]
   o, e := time.Parse(time.RFC3339, left + right)
   if e != nil {
      return 0, e
   }
   return time.Since(o).Hours(), nil
}

func Views(m winter.Map) (float64, error) {
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
