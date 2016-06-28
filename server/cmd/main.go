package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/oinume/lekcije/server/web"
)

var lekcijeEnv = os.Getenv("LEKCIJE_ENV")
var templateDir = "static"

func isProduction() bool {
	return lekcijeEnv == "production"
}

func templatePath() string {
	if isProduction() {
		return "static"
	} else {
		return "src/www"
	}
}

// TODO: move somewhere proper
var definedEnvs = map[string]string{
	"GOOGLE_CLIENT_ID":     "",
	"GOOGLE_CLIENT_SECRET": "",
}

func init() {
	if isProduction() {
		templateDir = "static"
	} else {
		templateDir = "src/www"
	}
	// Check env
	for key, _ := range definedEnvs {
		if value := os.Getenv(key); value != "" {
			definedEnvs[key] = value
		} else {
			log.Fatalf("Env %v is not defined.", key)
		}
	}
}

func main() {
	mux := mux()
	fmt.Println("Listening on :5000")
	http.ListenAndServe(":5000", mux)
}

func mux() *goji.Mux {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/"), index)
	mux.HandleFuncC(pat.Get("/oauth/google"), web.OAuthGoogle)
	//mux.HandleFuncC(pat.Get("/:name"), index)
	return mux
}

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(templateDir, "index.html")))
	if err := t.Execute(w, nil); err != nil {
		panic(err)
	}
}

//func loadTemplate()
//func index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	name := pat.Param(ctx, "name")
//	if name == "" {
//		name = r.URL.Query().Get("name")
//	}
//	fmt.Fprintf(w, "Hello, %s!", name)
//}
