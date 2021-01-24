package main

import (
   "encoding/json"
   "github.com/89z/winter/youtube"
   "log"
   "net/url"
   "os"
   "regexp"
   "strconv"
   "strings"
   "time"
)

func main() {
   if len(os.Args) != 2 {
      println("youtube-insert <URL>")
      os.Exit(1)
   }
   url_s := os.Args[1]
   o, e := url.Parse(url_s)
   check(e)
   id := o.Query().Get("v")
   // year
   info_m, e := youtube.Info(id)
   check(e)
   if info_m["description"] == nil {
      log.Fatal("Clapham Junction")
   }
   desc := info_m.M("description").S("simpleText")
   year_s := info_m.S("publishDate")[:4]
   /* the order doesnt matter here, as we will find the lowest date of all
   matches */
   reg_a := []string{
      ` (\d{4})`, `(\d{4}) `, `Released on: (\d{4})`, `℗ (\d{4})`,
   }
   for _, reg_s := range reg_a {
      mat_s := findSubmatch(reg_s, desc)
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
   line := regexp.MustCompile(".* · .*").FindString(desc)
   if line != "" {
      title_a := strings.Split(line, " · ")
      artist_a := title_a[1:]
      title_s = strings.Join(artist_a, ", ") + " - " + title_a[0]
   }
   // time
   date_n := time.Now().Unix()
   date_s := strconv.FormatInt(date_n, 36)
   // image
   image_s := getImage(id)
   // print
   rec_a := []interface{}{date_s, year_n, "y/" + id + image_s, title_s}
   rec, e := json.Marshal(rec_a)
   check(e)
   rec = append(rec, ',', '\n')
   os.Stdout.Write(rec)
}
