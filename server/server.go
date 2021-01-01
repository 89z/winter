package main

import (
   "database/sql"
   "html/template"
   "log"
   "net/http"
   "net/url"
   "time"
   _ "github.com/mattn/go-sqlite3"
)

func check(err error) {
   if err != nil {
      log.Fatal(err)
   }
}

type Page struct {
   H1 url.Values
   H2 time.Time
   Table []string
}

const tpl = `
<h1>{{ .H1 }}</h1>
<h2>{{ .H2 }}</h2>
{{ range .Table }}
<div>{{ . }}</div>
{{ end }}
`

func SelectOne(w http.ResponseWriter, r *http.Request) {
   winter_s := os.Getenv("WINTER")
   db, e := sql.Open("sqlite3", winter_s)
   check(e)
   query_o, e := db.Query("select artist_s from artist_t")
   check(e)
   var artist_a []string
   for query_o.Next() {
      var artist string
      e = query_o.Scan(&artist)
      check(e)
      artist_a = append(artist_a, artist)
   }
   t, e := template.New("webpage").Parse(tpl)
   check(e)
   data := Page{
      H1: r.URL.Query(),
      H2: time.Now(),
      Table: artist_a,
   }
   e = t.Execute(w, data)
   check(e)
}

func main() {
   http.HandleFunc("/", SelectOne)
   srv := http.Server{}
   srv.ListenAndServe()
}
