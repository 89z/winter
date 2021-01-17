package main

import (
   "encoding/json"
   "net/url"
   "os"
   "regexp"
   "strconv"
   "strings"
   "time"
   "winter"
   "winter/youtube"
)

func main() {
   if len(os.Args) != 2 {
      println("youtube-insert <URL>")
      os.Exit(1)
   }
   url := os.Args[1]
   o, e := url.Parse(url)
   check(e)
   id := o.Query().Get("v")
   // year
   info_m, e := youtube.Info(id)
   check(e)
   if info_m["description"] == nil {
      log.Fatal("Clapham Junction")
   }
   desc_s := info_m.M("description").S("simpleText")
   year_s := info_m.S("publishDate")[:4]
   /* the order doesnt matter here, as we will find the lowest date of all
   matches */
   reg_a := []string{
      ` (\d{4})`, `(\d{4}) `, `Released on: (\d{4})`, `℗ (\d{4})`,
   }
   for _, reg_s := range reg_a {
      mat_s := FindSubmatch(reg_s, desc_s)
      if mat_s == "" {
         continue
      }
      if mat_s >= year_s {
         continue
      }
      year_s = mat_s
   }
   year_n, e := strconv.Atoi(year_s)
   check(e)
   // song, artist
   title_s := info_m.M("title").S("simpleText")
   line_s := regexp.MustCompile(".* · .*").FindString(desc_s)
   if line_s != "" {
      title_a := strings.Split(line_s, " · ")
      artist_a := title_a[1:]
      title_s = strings.Join(artist_a, ", ") + " - " + title_a[0]
   }
   // time
   date_n := time.Now().Unix()
   date_s := strconv.FormatInt(date_n, 36)
   // image
   image_s := GetImage(id)
   // print
   rec_a := winter.Slice{date_s, year_n, "y/" + id + image_s, title_s}
   json_y, e := json.Marshal(rec_a)
   check(e)
   json_y = append(json_y, ',', '\n')
   os.Stdout.Write(json_y)
}
