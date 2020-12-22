package main

func LocalAlbum() {
   foreach ($artist_o as $album_s => $o_album) {
      $good_b = false;
      $done_b = true;
      foreach ($o_album as $track_s => $rate_s) {
         if ($track_s == '@id') {
            $local_m[$album_s] = 'black';
            continue 2;
         }
         if ($rate_s == 'good') {
            $good_b = true;
         }
         if ($rate_s == '') {
            $done_b = false;
         }
      }
      if ($good_b && $done_b) {
         $local_m[$album_s] = 'green';
      }
      if ($good_b && ! $done_b) {
         $local_m[$album_s] = 'lightgreen';
      }
      if (! $good_b && $done_b) {
         $local_m[$album_s] = 'red';
      }
      if (! $good_b && ! $done_b) {
         $local_m[$album_s] = 'lightred';
      }
   }
   return $local_m;
}

# remote albums
function mb_albums(string $arid_s): array {
   $query_m['artist'] = $arid_s;
   $query_m['fmt'] = 'json';
   $query_m['inc'] = 'release-groups';
   $query_m['limit'] = 100;
   $query_m['offset'] = 0;
   $query_m['status'] = 'official';
   $query_m['type'] = 'album';
   $remote_m = [];
   $url_r = curl_init();
   curl_setopt($url_r, CURLOPT_RETURNTRANSFER, true);
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
