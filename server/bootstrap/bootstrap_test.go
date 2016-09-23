package bootstrap

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCheckEnvs(t *testing.T) {
	a := assert.New(t)
	CheckEnvs()
	a.NotEmpty(Envs.DBURL)
}
