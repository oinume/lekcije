package main

import (
    "net/http"
    "fmt"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", root)
    http.ListenAndServe(":5000", mux)
}

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
