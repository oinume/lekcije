package http

import (
	"net/http"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/model"
)

// GET /api/status
func (s *server) getAPIStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]bool{
		"db": true,
	}

	db, err := model.OpenDB(config.DefaultVars.DBURL(), 1, config.DefaultVars.DebugSQL)
	if err == nil {
		defer db.Close()
		if err := db.DB().Ping(); err != nil {
			data["db"] = false
		}
	} else {
		data["db"] = false
	}

	for _, status := range data {
		if !status {
			writeJSON(w, http.StatusInternalServerError, data)
			return
		}
	}
	writeJSON(w, http.StatusOK, data)
}
