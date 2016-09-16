package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/model"
	"github.com/pkg/errors"
)

// GET /api/status
func GetStatus(w http.ResponseWriter, r *http.Request) {
	db, err := model.Open(os.Getenv("DB_DSN"))
	if err != nil {
		controller.InternalServerError(w, fmt.Errorf("Failed to model.Open(): err=%v", err))
		return
	}
	if err := db.DB().Ping(); err != nil {
		controller.InternalServerError(w, fmt.Errorf("Failed to DB.Ping(): err=%v", err))
		return
	}
	data := map[string]bool{
		"db": true,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		controller.InternalServerError(w, fmt.Errorf("Failed to encode JSON: err=%v", err))
		return
	}
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
