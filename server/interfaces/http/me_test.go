package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMe(t *testing.T) {
	//t.Parallel()
	a := assert.New(t)
	r := require.New(t)
	s := NewServer(&interfaces.ServerArgs{
		DB: helper.DB(t),
	})

	req, err := http.NewRequest("GET", "/me", nil)
	r.NoError(err)
	user := helper.CreateRandomUser(t)
	ctx := context_data.SetLoggedInUser(req.Context(), user)
	ctx = context_data.SetTrackingID(ctx, fmt.Sprint(user.ID))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	http.HandlerFunc(s.getMeHandler()).ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)
	a.Contains(w.Body.String(), "フォローしている講師")
}

func TestGetMeSetting(t *testing.T) {
	//t.Parallel()
	a := assert.New(t)
	r := require.New(t)
	s := NewServer(&interfaces.ServerArgs{
		DB: helper.DB(t),
	})

	req, err := http.NewRequest("GET", "/me/setting", nil)
	r.NoError(err)
	user := helper.CreateRandomUser(t)
	ctx := context_data.SetLoggedInUser(req.Context(), user)
	ctx = context_data.SetTrackingID(ctx, fmt.Sprint(user.ID))
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	http.HandlerFunc(s.getMeHandler()).ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)
	a.Contains(w.Body.String(), "フォローしている講師")
}
