package main

import (
   "log"
   "net/url"
   "os"
   "winter/youtube"
)

func main() {
   if len(os.Args) != 2 {
      println("youtube-track <URL>")
      os.Exit(1)
   }
   url_s := os.Args[1]
   o, e := url.Parse(url_s)
   if e != nil {
      log.Fatal(e)
   }
   id_s := o.Query().Get("v")
   // year
   m, e := youtube.Info(id_s)
   if e != nil {
      log.Fatal(e)
   }
   if m["description"] == nil {
      log.Fatal("Clapham Junction")
   }
   desc_s := m.M("description").S("simpleText")
   year_s := m.S("publishDate")
   reg_a := []string{
      ` (\d{4})`, `(\d{4,}) `, `Released on: (\d{4})`, `℗ (\d{4})`,
   }
   for _, reg_s := range reg_a {
      $mat_n = preg_match($reg_s, $info_o->description->simpleText, $mat_a);
      if ($mat_n === 0) {
         continue;
      }
      $mat_s = $mat_a[1];
      if ($mat_s >= $year_s) {
         continue;
      }
      $year_s = $mat_s;
   }
   $year_n = (int)($year_s);
   # song, artist
   $mat_n = preg_match('/.* · .*/', $info_o->description->simpleText, $line_a);
   if ($mat_n !== 0) {
      $line_s = $line_a[0];
      $title_a = explode(' · ', $line_s);
      $artist_a = array_slice($title_a, 1);
      $title_s = implode(', ', $artist_a) . ' - ' . $title_a[0];
   } else {
      $title_s = $info_o->title->simpleText;
   }
   # time
   function encode36(int $n): string {
      $s = (string) $n;
      return base_convert($s, 10, 36);
   }
   $date_n = time();
   $date_s = encode36($date_n);
   # image
   $jpg_a = ['/sddefault','/sd1', '/hqdefault'];
   foreach ($jpg_a as $jpg_s) {
      $url_s = 'https://i.ytimg.com/vi/' . $id_s . $jpg_s . '.jpg';
      echo $url_s, "\n";
      $head_a = get_headers($url_s);
      $code_s = $head_a[0];
      if (str_contains($code_s, '200 OK')) {
         break;
      }
   }
   if ($jpg_s == '/sddefault') {
      $jpg_s = '';
   }
   # print
   $rec_a = [$date_s, $year_n, 'y/' . $id_s . $jpg_s, $title_s];
   $json_s = json_encode($rec_a, JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
   echo $json_s, ",\n";
}
