package main

import (
   "html/template"
   "log"
   "net/http"
   "net/url"
   "time"
)

func check(err error) {
   if err != nil {
      log.Fatal(err)
   }
}

type Page struct {
   Items url.Values
   Title time.Time
}

const tpl = `
<h1>{{ .Title }}</h1>
{{ range .Items }}
<div>{{ . }}</div>
{{ end }}
`

func Index(w http.ResponseWriter, r *http.Request) {
   t, e := template.New("webpage").Parse(tpl)
   check(e)
   data := Page{
      Items: r.URL.Query(),
      Title: time.Now(),
   }
   e = t.Execute(w, data)
   check(e)
}

func main() {
   http.HandleFunc("/", Index)
   srv := http.Server{}
   srv.ListenAndServe()
}
