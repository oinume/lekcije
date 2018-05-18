package http

import (
	"net/http"
	"os"
	"strings"
)

// GET /api/debug/envVar
func GetAPIDebugEnvVar(w http.ResponseWriter, r *http.Request) {
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
		JSON(w, http.StatusOK, vars)
	} else {
		JSON(w, http.StatusOK, struct{}{})
	}
}

// GET /api/debug/httpHeader
func GetAPIDebugHTTPHeader(w http.ResponseWriter, r *http.Request) {
	debugKey := r.FormValue("debugKey")
	envDebugKey := os.Getenv("DEBUG_KEY")
	if debugKey != "" && debugKey == envDebugKey {
		headers := make(map[string]string)
		for name := range r.Header {
			headers[name] = r.Header.Get(name)
		}
		JSON(w, http.StatusOK, headers)
	} else {
		JSON(w, http.StatusOK, struct{}{})
	}
}
