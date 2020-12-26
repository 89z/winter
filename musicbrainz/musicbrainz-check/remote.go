package main

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "winter/snow"
)

type Remote struct {
   Date string
   Title string
}

var (
   offset_n float64
   remote_a []Remote
)

func RemoteAlbum(mb_s string) ([]Remote, error) {
   q := url.Values{}
   q.Set("artist", mb_s)
   q.Set("fmt", "json")
   q.Set("inc", "release-groups")
   q.Set("limit", "100")
   q.Set("status", "official")
   q.Set("type", "album")
   for {
      url_s := "https://musicbrainz.org/ws/2/release?" + q.Encode()
      fmt.Println(url_s)
      o, e := http.Get(url_s)
      if e != nil {
         return nil, e
      }
      json_m := snow.Map{}
      e = json.NewDecoder(o.Body).Decode(&json_m)
      if e != nil {
         return nil, e
      }
      release_a := json_m.A("releases")
      for n := range release_a {
         group_m := release_a.M(n).M("release-group")
         second_a := group_m.A("secondary-types")
         if len(second_a) > 0 {
            continue
         }
         /* right here we want to take the Release title, not the Release Group
         title. Two reasons. First, some releases in the same group are quite
         different, and we may want to get both of them:
         musicbrainz.org/release-group/ec1842b2-9393-3779-aa80-280b90a55bef

         second, some releases have slight name difference. For example, at this
         time, under this release group:
         musicbrainz.org/release-group/58d91845-3734-344b-ba10-fcae524f22c1

         the group title is using U+2010 HYPHEN, while one of the releases is
         using U+002D HYPHEN-MINUS. So in the case of different album content,
         we will be made aware of the different release choices and can get both
         and just mark duplicate songs as needed. In the case of spelling
         differences, we can make edits to MusicBrainz website, then re run the
         check */
         album_s := group_m.S("title")
         remote_a = append(remote_a, Remote{
            group_m.S("first-release-date"), album_s,
         })
      }
      offset_n += 100
      if offset_n >= json_m.N("release-count") {
         break
      }
      q.Set("offset", fmt.Sprint(offset_n))
   }
   return remote_a, nil
}
