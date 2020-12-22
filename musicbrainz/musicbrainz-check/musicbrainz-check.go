package main
import "os"

func main() {
   if len(os.Args) != 2 {
      println("musicbrainz-check <artist>")
      os.Exit(1)
   }
   artist_s := os.Args[1]
   // local albums
   $arid_s = $local_o->$artist_s->{'@mb'};
   $local_m = si_color($local_o->$artist_s);
   $remote_m = mb_albums($arid_s);
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
