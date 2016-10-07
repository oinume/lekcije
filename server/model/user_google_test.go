package model

import (
	"testing"

	"github.com/oinume/lekcije/server/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserGoogleService_FindByPK(t *testing.T) {
	a := assert.New(t)

	created, err := userGoogleService.Create("1", 1)
	a.NoError(err)
	a.Equal("1", created.GoogleID)

	_, err = userGoogleService.FindByPK("0")
	a.IsType(&errors.NotFound{}, err)
}
