package http

import (
	"encoding/json"
	"net/http"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model"
)

// GET /api/status
func (s *server) getAPIStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Include connection statistics
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

// GET /api/me/followingTeachers
func (s *server) getAPIMeFollowingTeachers(w http.ResponseWriter, r *http.Request) {
	// SELECT t.id, t.name FROM following_teachers AS ft
	// INNER JOIN teachers AS t ON ft.teacher_id = t.id
	// WHERE ft.user_id = ?
	// ORDER BY ft.updated_at
	teachers := []map[string]interface{}{
		{"id": "1", "name": "Xai"},
		{"id": "2", "name": "Emina"},
		{"id": "3", "name": "Tasha"},
	}
	if err := json.NewEncoder(w).Encode(teachers); err != nil {
		// TODO: JSON
		internalServerError(s.appLogger, w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to encode JSON"),
		), 0)
		return
	}
}
