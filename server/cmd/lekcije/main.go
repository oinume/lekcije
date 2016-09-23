package main

import (
	"fmt"
	"net/http"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/mux"
)

func init() {
	bootstrap.CheckEnvs()
}

func main() {
	port := config.ListenPort()
	mux := mux.Create()
	fmt.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
