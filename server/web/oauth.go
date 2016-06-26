package web

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint:     google.Endpoint,
	RedirectURL:  fmt.Sprintf("http://127.0.0.1:%d/callback", 5000),
	Scopes: []string{
		"openid email",
	},
}

func OAuthGoogle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, googleOAuthConfig.AuthCodeURL(""), http.StatusFound)
}
