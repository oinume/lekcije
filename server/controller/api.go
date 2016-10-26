package controller

import (
	"encoding/json"
	"net/http"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

// GET /api/status
func GetAPIStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Include connection statistics
	data := map[string]bool{
		"db":    true,
		"redis": true,
	}

	db, err := model.OpenDB(bootstrap.HTTPServerEnvVars.DBURL)
	if err == nil {
		defer db.Close()
		if err := db.DB().Ping(); err != nil {
			data["db"] = false
		}
	} else {
		data["db"] = false
	}

	redis, err := model.OpenRedis(bootstrap.HTTPServerEnvVars.RedisURL)
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
			JSON(w, http.StatusInternalServerError, data)
			return
		}
	}
	JSON(w, http.StatusOK, data)
}

// GET /api/me/followingTeachers
func GetAPIMeFollowingTeachers(w http.ResponseWriter, r *http.Request) {
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
		InternalServerError(w, errors.InternalWrapf(err, "Failed to encode JSON"))
		return
	}
}
