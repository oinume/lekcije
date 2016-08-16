package e2e

import (
	"testing"

	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	a := assert.New(t)
	driver := agouti.ChromeDriver()
	//driver := agouti.PhantomJS()
	driver.HTTPClient = client
	err := driver.Start()
	a.NoError(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	a.NoError(err)
	a.NoError(page.Navigate(server.URL))
	//time.Sleep(10 * time.Second)
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
