package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/model"
	"github.com/pkg/errors"
)

// GET /api/status
func GetStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]bool{
		"db":    true,
		"redis": true,
	}

	db, err := model.OpenDB(os.Getenv("DB_DSN"))
	if err == nil {
		if err := db.DB().Ping(); err != nil {
			data["db"] = false
		}
	} else {
		data["db"] = false
	}

	redis, err := model.OpenRedis(os.Getenv("REDIS_URL"))
	if err == nil {
		if redis.Ping().Err() != nil {
			data["redis"] = false
		}
	} else {
		data["redis"] = false
	}

	for _, status := range data {
		if !status {
			controller.JSON(w, http.StatusInternalServerError, data)
			return
		}
	}
	controller.JSON(w, http.StatusOK, data)
}

// GET /api/me/followingTeachers
func GetMeFollowingTeachers(w http.ResponseWriter, r *http.Request) {
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
		controller.InternalServerError(w, errors.Wrapf(err, "Failed to encode JSON"))
		return
	}
}
