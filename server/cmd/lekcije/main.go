package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/mux"
)

// TODO: move somewhere proper and make it be struct
var definedEnvs = map[string]string{
	"GOOGLE_CLIENT_ID":     "",
	"GOOGLE_CLIENT_SECRET": "",
	"DB_URL":               "",
	"NODE_ENV":             "",
	"ENCRYPTION_KEY":       "",
}

func init() {
	// Check env
	for key, _ := range definedEnvs {
		if value := os.Getenv(key); value != "" {
			definedEnvs[key] = value
		} else {
			log.Fatalf("Env '%v' must be defined.", key)
		}
	}

	http.DefaultClient.Timeout = 5 * time.Second
}

func main() {
	port := config.ListenPort()
	mux := mux.Create()
	fmt.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
