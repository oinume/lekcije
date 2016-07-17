package web

import (
	"html/template"
	"net/http"
	"path"

	"golang.org/x/net/context"
)

func Root(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join(TemplateDir, "index.html")))
	if err := t.Execute(w, nil); err != nil {
		panic(err)
	}
}
