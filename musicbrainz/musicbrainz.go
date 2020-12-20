package release

import (
   "encoding/json"
   "net/http"
   "net/url"
)

type MB struct {
   API string
   ID string
   Query url.Values
}

func NewDecode(mbid_s string) MB {
   return MB{
      "https://musicbrainz.org/ws/2/release",
      mbid_s,
      url.Values{
         "fmt": []string{"json"},
         "inc": []string{"artist-credits recordings"},
      },
   }
}

func (mb_o MB) Group() (Slice, error) {
   mb_o.Query["release-group"] = []string{mb_o.ID}
   url_s := mb_o.API + "?" + mb_o.Query.Encode()
   println(url_s)
   get_o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   m := Map{}
   e = json.NewDecoder(get_o.Body).Decode(&m)
   if e != nil {
      return nil, e
   }
   return m.A("releases"), nil
}

func (mb_o MB) Release() (Map, error) {
   url_s := mb_o.API + "/" + mb_o.ID + "?" + mb_o.Query.Encode()
   println(url_s)
   get_o, e := http.Get(url_s)
   if e != nil {
      return nil, e
   }
   m := Map{}
   e = json.NewDecoder(get_o.Body).Decode(&m)
   if e != nil {
      return nil, e
   }
   return m, nil
}
