package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/errors"
)

var _ = fmt.Print

func (s *server) static(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func (s *server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := context_data.GetLoggedInUser(r.Context()); err == nil {
			http.Redirect(w, r, "/me", http.StatusFound)
		} else {
			s.indexLogout(w, r)
		}
	}
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	if _, err := context_data.GetLoggedInUser(r.Context()); err == nil {
		http.Redirect(w, r, "/me", http.StatusFound)
	} else {
		s.indexLogout(w, r)
	}
}

func (s *server) indexLogout(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("index.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, false, 0),
	}

	if err := t.Execute(w, data); err != nil {
		internalServerError(s.appLogger, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), 0)
		return
	}
}

func (s *server) signupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.signup(w, r)
	}
}

func (s *server) signup(w http.ResponseWriter, r *http.Request) {
	if _, err := context_data.GetLoggedInUser(r.Context()); err == nil {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}

	t := ParseHTMLTemplates(TemplatePath("signup.html"))
	data := struct {
		commonTemplateData
	}{
		commonTemplateData: s.getCommonTemplateData(r, false, 0),
	}
	if err := t.Execute(w, &data); err != nil {
		internalServerError(s.appLogger, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), 0)
		return
	}
}

func (s *server) robotsTxtHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.robotsTxt(w, r)
	}
}

func (s *server) robotsTxt(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf(`
User-agent: *
Allow: /
Sitemap: %s/sitemap.xml
`, config.WebURL())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, strings.TrimSpace(content))
}

func (s *server) sitemapXMLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.sitemapXML(w, r)
	}
}

func (s *server) sitemapXML(w http.ResponseWriter, r *http.Request) {
	// TODO: Move to an external file
	content := fmt.Sprintf(`
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>%s/</loc>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>%s/signup</loc>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>%s/terms</loc>
    <priority>1.0</priority>
  </url>
</urlset>
	`, config.WebURL(), config.WebURL(), config.WebURL())
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, strings.TrimSpace(content))
}

func (s *server) termsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.terms(w, r)
	}
}

func (s *server) terms(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("terms.html"))
	type Data struct {
		commonTemplateData
	}

	userID := uint32(0)
	if user, err := context_data.GetLoggedInUser(r.Context()); err == nil {
		userID = user.ID
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, false, userID),
	}

	if err := t.Execute(w, data); err != nil {
		internalServerError(s.appLogger, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to template.Execute()"),
		), 0)
		return
	}
}
