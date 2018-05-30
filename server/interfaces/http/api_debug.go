package http

import (
	"net/http"
	"os"
	"strings"
)

func (s *server) getAPIDebugEnvVarHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getAPIDebugEnvVar(w, r)
	}
}

// GET /api/debug/envVar
func (s *server) getAPIDebugEnvVar(w http.ResponseWriter, r *http.Request) {
	debugKey := r.FormValue("debugKey")
	envDebugKey := os.Getenv("DEBUG_KEY")
	if debugKey != "" && debugKey == envDebugKey {
		vars := make(map[string]string)
		for _, v := range os.Environ() {
			kv := strings.SplitN(v, "=", 2)
			if kv[0] == "ENCRYPTION_KEY" || kv[0] == "DEBUG_KEY" {
				continue
			}
			vars[kv[0]] = kv[1]
		}
		writeJSON(w, http.StatusOK, vars)
	} else {
		writeJSON(w, http.StatusOK, struct{}{})
	}
}

func (s *server) getAPIDebugHTTPHeaderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getAPIDebugHTTPHeader(w, r)
	}
}

// GET /api/debug/httpHeader
func (s *server) getAPIDebugHTTPHeader(w http.ResponseWriter, r *http.Request) {
	debugKey := r.FormValue("debugKey")
	envDebugKey := os.Getenv("DEBUG_KEY")
	if debugKey != "" && debugKey == envDebugKey {
		headers := make(map[string]string)
		for name := range r.Header {
			headers[name] = r.Header.Get(name)
		}
		writeJSON(w, http.StatusOK, headers)
	} else {
		writeJSON(w, http.StatusOK, struct{}{})
	}
}
