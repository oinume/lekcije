package web

// TODO: Create package 'api'

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// /api/me/followingTeachers
func ApiGetMeFollowingTeachers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		InternalServerError(w, errors.Wrapf(err, "Failed to encode JSON"))
		return
	}
}
