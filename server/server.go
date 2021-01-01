package main

import (
   "fmt"
   "net/http"
   "time"
)

func Index(w http.ResponseWriter, r *http.Request) {
   t := time.Now()
   q := r.URL.Query()
   fmt.Fprint(w, t, "\n", q, "\n")
}

func Month(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintln(w, `
   <!doctype html>
   <html lang="en">
   <head>
      <meta charset="utf-8">
      <title>March</title>
   </head>
   <body>
      <h1>March</h1>
   </body>
   </html>
   `)
}

func main() {
   http.HandleFunc("/", Index)
   http.HandleFunc("/month", Month)
   srv := http.Server{}
   srv.ListenAndServe()
}
