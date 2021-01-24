package musicbrainz

import (
   "github.com/89z/x"
   "github.com/89z/x/json"
   "net/url"
   "path"
   "strings"
)

func Release(in string) (x.Map, error) {
   id := path.Base(in)
   v := url.Values{}
   v.Set("fmt", "json")
   v.Set("inc", "artist-credits recordings")
   if strings.Contains(in, "/release/") {
      out := "https://musicbrainz.org/ws/2/release/" + id + "?" + v.Encode()
      return json.LoadHttp(out)
   }
   v.Set("release-group", id)
   out := "https://musicbrainz.org/ws/2/release?" + v.Encode()
   group, e := json.LoadHttp(out)
   if e != nil {
      return nil, e
   }
   albums := group.A("releases")
   Sort(albums)
   return albums.M(0), nil
}
