package main

import (
   "database/sql"
   "net/url"
   "winter/snow"
)

func LocalAlbum(open_o *sql.DB, artist_s string) (map[string]string, error) {
   query_o, e := open_o.Query(`
   select
      album_s,
      url_s,
      count(1) filter (where note_s = '') as unrated_n,
      count(1) filter (where note_s = 'good') as good_n
   from album_t
   natural join song_album_t
   natural join song_t
   natural join song_artist_t
   natural join artist_t
   where artist_s = ?
   group by album_n
   `, artist_s)
   if e != nil {
      return nil, e
   }
   var (
      album_s string
      url_s string
      unrated_n int
      good_n int
   )
   local_m := map[string]string{}
   for query_o.Next() {
      e = query_o.Scan(&album_s, &url_s, &unrated_n, &good_n)
      if e != nil {
         return nil, e
      }
      if snow.Pop(url_s) {
         local_m[album_s] = "green"
         continue
      }
      /*
      unrated tracks | good tracks | color
      ---------------|-------------|------
      0              | 0           | red
      0              | 1           | green
      1              | 0           | light red
      1              | 1           | light green
      */
      if unrated_n == 0 && good_n == 0 {
         local_m[album_s] = "red"
         continue
      }
      if unrated_n == 0 {
         local_m[album_s] = "green"
         continue
      }
      if good_n == 0 {
         local_m[album_s] = "light red"
         continue
      }
      local_m[album_s] = "light green"
   }
   return local_m, nil
}

function RemoteAlbum(mb_s string) map[string]string {
   q := url.Values{}
   q.Set("artist", mb_s)
   q.Set("fmt", "json")
   q.Set("inc", "release-groups")
   q.Set("limit", "100")
   q.Set("offset", "0")
   q.Set("status", "official")
   q.Set("type", "album")
   remote_m := map[string]string{}
   curl_setopt($url_r, CURLOPT_USERAGENT, 'anonymous');
   while (true) {
      # part 1
      $query_s = http_build_query($query_m);
      $url_s = 'https://musicbrainz.org/ws/2/release?' . $query_s;
      curl_setopt($url_r, CURLOPT_URL, $url_s);
      echo $url_s, "\n";
      # part 2
      $json_s = curl_exec($url_r);
      # part 3
      $remote_o = json_decode($json_s);
      foreach ($remote_o->releases as $o_re) {
         $o_rg = $o_re->{'release-group'};
         $a_sec = $o_rg->{'secondary-types'};
         if (count($a_sec) > 0) {
            continue;
         }
         if (array_key_exists($o_rg->title, $remote_m)) {
            continue;
         }
         $remote_m[$o_rg->title] = $o_rg->{'first-release-date'};
      }
      $query_m['offset'] += $query_m['limit'];
      if ($query_m['offset'] >= $remote_o->{'release-count'}) {
         break;
      }
   }
   return $remote_m;
}
