package main
import "net/http"

func main() {
   http.HandleFunc("/", selectOne)
   serve := http.Server{}
   serve.ListenAndServe()
}
