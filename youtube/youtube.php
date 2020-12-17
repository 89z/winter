<?php
declare(strict_types = 1);
extension_loaded('openssl') or die('openssl');
require_once 'cove/color-string.php';

function format_number(float $n): string {
   $n2 = (int)(log10($n) / 3);
   return sprintf('%.3f', $n / 1e3 ** $n2) . ['', ' K', ' M', ' B'][$n2];
}

function seconds(object $new_o, object $old_o): float {
   $o = $new_o->diff($old_o);
   return $o->days * 86400 + $o->h * 3600 + $o->i * 60 + $o->s + $o->f;
}

function youtube_info(string $id_s): object {
   # part 1
   $info_s = 'https://www.youtube.com/get_video_info?video_id=' . $id_s;
   echo $info_s, "\n";
   # part 2
   $get_s = file_get_contents($info_s);
   parse_str($get_s, $get_m);
   # part 3
   $resp_s = $get_m['player_response'];
   # part 4
   return json_decode($resp_s)->microformat->playerMicroformatRenderer;
}

function youtube_views(object $info_o): string {
   $views_n = (int)($info_o->viewCount);
   $date_o = DateTime::createFromFormat('!Y-m-d', $info_o->publishDate);
   $sec_n = seconds(new DateTime, $date_o);
   $rate_n = $views_n / ($sec_n / 60 / 60 / 24 / 365);
   $rate_s = format_number($rate_n);
   if ($rate_n > 8_000_000) {
      return color_red($rate_s);
   }
   return color_green($rate_s);
}
