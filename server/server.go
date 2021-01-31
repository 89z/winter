package main

import (
   "database/sql"
   "github.com/89z/x"
   "html/template"
   "net/http"
   "net/url"
   "os"
   "time"
   _ "github.com/mattn/go-sqlite3"
)

const tpl = `
<h1>{{ .H1 }}</h1>
<h2>{{ .H2 }}</h2>
{{ range .Table }}
<div>{{ . }}</div>
{{ end }}
`

func selectOne(w http.ResponseWriter, r *http.Request) {
   db, e := sql.Open(
      "sqlite3", os.Getenv("WINTER"),
   )
   x.Check(e)
   query, e := db.Query("select artist_s from artist_t")
   x.Check(e)
   var artists []string
   for query.Next() {
      var artist string
      e = query.Scan(&artist)
      x.Check(e)
      artists = append(artists, artist)
   }
   t, e := template.New("webpage").Parse(tpl)
   x.Check(e)
   data := page{
      h1: r.URL.Query(),
      h2: time.Now(),
      table: artists,
   }
   e = t.Execute(w, data)
   x.Check(e)
}

type page struct {
   h1 url.Values
   h2 time.Time
   table []string
}
