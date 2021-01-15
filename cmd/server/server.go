package main
import "net/http"

func main() {
   http.HandleFunc("/", SelectOne)
   serve := http.Server{}
   serve.ListenAndServe()
}
