package main

import (
	"fmt"
	"net/http"

	"html/template"
	"os"
	"path"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
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

func init() {
	if isProduction() {
		templateDir = "static"
	} else {
		templateDir = "src/www"
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
