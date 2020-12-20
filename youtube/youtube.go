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
   "winter/assert"
   "winter/color"
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

func timeHours(value string) (float64, error) {
   layout := time.RFC3339[:len(value)]
   o, e := time.Parse(layout, value)
   if e != nil {
      return 0, e
   }
   return time.Since(o).Hours(), nil
}

func Info(id_s string) (assert.Map, error) {
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
   json_m := assert.Map{}
   e = json.Unmarshal([]byte(resp_s), &json_m)
   if e != nil {
      return nil, e
   }
   return json_m.M("microformat").M("playerMicroformatRenderer"), nil
}

func Views(info_m assert.Map) (string, error) {
   view_s := info_m.S("viewCount")
   view_n, e := floatVal(view_s)
   if e != nil {
      return "", e
   }
   date_s := info_m.S("publishDate")
   hour_n, e := timeHours(date_s)
   if e != nil {
      return "", e
   }
   rate_n := view_n / (hour_n / 24 / 365)
   rate_s := numberFormat(rate_n)
   if rate_n > 8_000_000 {
      return color.Red(rate_s), nil
   }
   return color.Green(rate_s), nil
}
