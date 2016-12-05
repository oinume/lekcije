package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/mux"
)

func init() {
	bootstrap.CheckServerEnvVars()
}

func main() {
	port := config.ListenPort()
	mux := mux.Create()
	fmt.Printf("Listening on :%v\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), mux); err != nil {
		log.Fatalf("ListenAndServe() on :%v failed: err = %v", port, err)
	}
}
