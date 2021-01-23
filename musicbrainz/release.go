package musicbrainz

import (
   "github.com/89z/x"
   "github.com/89z/x/json"
   "net/url"
   "sort"
)

func Group(id string) (x.Slice, error) {
   q := url.Values{}
   q.Set("fmt", "json")
   q.Set("inc", "artist-credits recordings")
   q.Set("release-group", id)
   url := "https://musicbrainz.org/ws/2/release?" + q.Encode()
   group, e := json.LoadHttp(url)
   if e != nil {
      return nil, e
   }
   return group.A("releases"), nil
}

func Release(id string) (x.Map, error) {
   q := url.Values{}
   q.Set("fmt", "json")
   q.Set("inc", "artist-credits recordings")
   url := "https://musicbrainz.org/ws/2/release/" + id + "?" + q.Encode()
   return json.LoadHttp(url)
}
