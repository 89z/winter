package main

import (
   "database/sql"
   "log"
   "os"
   _ "github.com/mattn/go-sqlite3"
)

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist_s := os.Args[1]
   // local albums
   db_s := os.Getenv("WINTER")
   open_o, e := sql.Open("sqlite3", db_s)
   if e != nil {
      log.Fatal(e)
   }
   var mb_s string
   e = open_o.QueryRow(
      "select mb_s from artist_t where artist_s LIKE ?", artist_s,
   ).Scan(&mb_s)
   if e != nil {
      log.Fatal(e)
   }
   $local_m = si_color($local_o->$artist_s);
   $remote_m = mb_albums($mb_s);
   arsort($remote_m);
   foreach ($remote_m as $title_s => $date_s) {
      echo $date_s, "\t";
      if (array_key_exists($title_s, $local_m)) {
         $class_s = $local_m[$title_s];
         printf('<td style="background:%s">%s', $class_s, $title_s);
      } else {
         printf('<td>%s', $title_s);
      }
   }
}
