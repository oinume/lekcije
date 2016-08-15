package e2e

import (
	"testing"

	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestIndex(t *testing.T) {
	a := assert.New(t)
	driver := agouti.ChromeDriver()
	//driver := agouti.PhantomJS()
	//driver.HTTPClient = client
	err := driver.Start()
	a.NoError(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	a.NoError(err)
	err = page.Navigate(server.URL)
	a.NoError(err)
	time.Sleep(10 * time.Second)

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
