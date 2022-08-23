package graphqltest

import (
	"net/http"
)

type authTransport struct {
	token  string
	parent http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return t.parent.RoundTrip(req)
}
