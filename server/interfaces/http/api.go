package http

import (
	"encoding/json"
	"net/http"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

func (s *server) getAPIStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getAPIStatus(w, r)
	}
}

// GET /api/status
func (s *server) getAPIStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Include connection statistics
	data := map[string]bool{
		"db":    true,
		"redis": true,
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

	redis, err := model.OpenRedis(config.DefaultVars.RedisURL)
	if err == nil {
		defer redis.Close()
		if redis.Ping().Err() != nil {
			data["redis"] = false
		}
	} else {
		data["redis"] = false
	}

	for _, status := range data {
		if !status {
			writeJSON(w, http.StatusInternalServerError, data)
			return
		}
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *server) getAPIMeFollowingTeachersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.getAPIMeFollowingTeachers(w, r)
	}
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
		internalServerError(w, errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessage("Failed to encode JSON"),
		))
		return
	}
}
