package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ = time.UTC
var _ = fmt.Print

func TestIndex(t *testing.T) {
	a := assert.New(t)
	driver := newWebDriver()
	err := driver.Start()
	a.NoError(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	a.NoError(err)
	a.NoError(page.Navigate(server.URL))
	button := page.Find("#btn-primary")
	a.NotEmpty(button.String())
}
