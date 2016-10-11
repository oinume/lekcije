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
	//if err := page.Navigate(oauthURL); err != nil {
	//	return "", err
	//}
	//if err := page.Find("#Email").Fill(email); err != nil {
	//	return "", err
	//}
	//if err := page.Find("#gaia_loginform").Submit(); err != nil {
	//	return "", err
	//}
}
