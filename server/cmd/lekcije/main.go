package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/web"
	"github.com/oinume/lekcije/server/web/middleware"
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
	"DB_DSN":               "",
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
			log.Fatalf("Env '%v' must be defined.", key)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	mux := mux()
	fmt.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}

func mux() *goji.Mux {
	mux := goji.NewMux()
	mux.UseC(middleware.SetDbToContext)
	mux.HandleFuncC(pat.Get("/"), index)
	mux.HandleFuncC(pat.Get("/status"), status)
	mux.HandleFuncC(pat.Get("/oauth/google"), web.OAuthGoogle)
	mux.HandleFuncC(pat.Get("/oauth/google/callback"), web.OAuthGoogleCallback)
	mux.HandleFuncC(pat.Get("/api/me/followingTeachers"), web.ApiGetMeFollowingTeachers)
	return mux
}

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(templateDir, "index.html")))
	if err := t.Execute(w, nil); err != nil {
		panic(err)
	}
}

func status(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db, err := model.Open()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to model.Open(): err=%v", err), http.StatusInternalServerError)
		return
	}
	if err := db.DB().Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to DB.Ping(): err=%v", err), http.StatusInternalServerError)
		return
	}
	data := map[string]bool{
		"db": true,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON"), http.StatusInternalServerError)
		return
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
