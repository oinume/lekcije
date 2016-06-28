package fetcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	fetcher := NewTeacherLessonFetcher(nil)
	teacher, err := fetcher.Fetch(1523)
	assert.NoError(err)
	assert.Equal("Crizelle", teacher.Name)
}
