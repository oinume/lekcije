package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ = time.UTC
var _ = fmt.Print

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
	//time.Sleep(10 * time.Second)

	err = page.Find("#Email").Fill(os.Getenv("E2E_GOOGLE_ACCOUNT"))
	a.NoError(err)
	page.Find("#gaia_loginform").Submit()
	a.NoError(err)

	time.Sleep(time.Second * 1)
	page.Find("#Passwd").Fill(os.Getenv("E2E_GOOGLE_PASSWORD"))
	a.NoError(err)
	page.Find("#gaia_loginform").Submit()
	a.NoError(err)

	time.Sleep(time.Second * 3)
	err = page.Find("#submit_approve_access").Click()
	a.NoError(err)
	//time.Sleep(time.Second * 10)
	// TODO: Check content
	// user.email == E2E_GOOGLE_ACCOUNT

	// TODO: GetCookies() returns error
	/*
		$ ag 'unexpected response' vendor/github.com/sclevine/agouti         ✱
		vendor/github.com/sclevine/agouti/api/internal/bus/client.go
		32:                     return fmt.Errorf("unexpected response: %s", responseBody)
	*/
	//cookies, err := page.GetCookies()
	//fmt.Printf("cookies = %+v, err = %v\n", cookies)
}

// TODO: user_api_token will be deleted after logout
