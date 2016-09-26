package bootstrap

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCheckCommonEnvVars(t *testing.T) {
	a := assert.New(t)
	CheckCLIEnvVars()
	a.NotEmpty(CLIEnvVars.DBURL)
}
