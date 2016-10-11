package e2e

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestOAuthGoogle(t *testing.T) {
	a := assert.New(t)
	driver := newWebDriver()
	err := driver.Start()
	a.NoError(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	a.NoError(err)
	a.NoError(page.Navigate(server.URL))
	//time.Sleep(10 * time.Second)
	link := page.FindByXPath("//div[@class='starter-template']/a")
	u, err := link.Attribute("href")
	a.NoError(err)
	fmt.Printf("u = %v, err = %v\n", u, err)
	link.Click()
}
