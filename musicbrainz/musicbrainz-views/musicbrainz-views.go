package main

import (
   "net/http"
   "net/url"
)

func GetContents(s string) ([]byte, error) {
   o, e := http.Get(s)
   if e != nil {
      return []byte{}, e
   }
   return ioutil.ReadAll(o.Body)
}

func FindSubmatch(pat string, sub []byte) []byte {
   a := regexp.MustCompile(pat).FindSubmatch(sub)
   if len(a) < 2 {
      return []byte{}
   }
   return a[1]
}

func YoutubeResult(query_s string) (string, error) {
   m := url.Values{}
   m.Set("search_query", query_s)
   res_s := "https://www.youtube.com/results?" + m.Encode()
   println(res_s)
   get_y, e := GetContents(res_s)
   if e != nil {
      return nil, e
   }
   find_y := FindSubmatch("/vi/([^/]*)/", get_y)
   return string(find_y), nil
}

func main() {
   if ($argc != 2) {
      echo <<<eof
usage:
musicbrainz-views.php <URL>

examples:
https://musicbrainz.org/release-group/d03bb6b1-d7b4-38ea-974e-847cbb31dca4
https://musicbrainz.org/release/7a629d52-6a61-3ea1-a0a0-dd50bdef63b4
eof;
      exit(1);
   }
   $url_s = $argv[1];
   $mbid_s = basename($url_s);
   $dec_o = new MusicBrainzDecode($mbid_s);
   if (str_contains($url_s, 'release-group')) {
      $rel_a = $dec_o->group();
      $rel_n = 0;
      foreach ($rel_a as $idx_n => $cur_o) {
         $rel_n = MusicBrainzReduce($rel_n, $cur_o, $idx_n, $rel_a);
      }
      $rel_o = $rel_a[$rel_n];
      echo 'musicbrainz.org/release/' . $rel_o->id, "\n";
   } else {
      $rel_o = $dec_o->release();
   }
   foreach ($rel_o->{'artist-credit'} as $artist_o) {
      $out_a[] = $artist_o->name;
   }
   $artists_s = implode(' ', $out_a);
   foreach ($rel_o->media as $media_o) {
      foreach ($media_o->tracks as $track_o) {
         $id_s = yt_result($artists_s . ' ' . $track_o->title);
         $info_o = youtube_info($id_s);
         echo youtube_views($info_o), "\n\n";
         usleep(500_000);
      }
   }
}
