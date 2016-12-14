package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"net/url"

	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
)

var _ = time.UTC
var _ = fmt.Print

func TestOAuthGoogleLogin(t *testing.T) {
	if os.Getenv("CIRCLECI") != "" {
		t.Skipf("Skip because it can't render Google login page.")
	}
	a := assert.New(t)
	driver := newWebDriver()
	err := driver.Start()
	a.Nil(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	a.Nil(err)
	a.Nil(page.Navigate(server.URL))
	//time.Sleep(10 * time.Second)
	link := page.FindByXPath("//div[@class='starter-template']/a")
	u, err := link.Attribute("href")
	a.Nil(err)
	fmt.Printf("u = %v, err = %v\n", u, err)
	link.Click()
	//time.Sleep(10 * time.Second)

	googleAccount := os.Getenv("E2E_GOOGLE_ACCOUNT")
	err = page.Find("#Email").Fill(googleAccount)
	a.Nil(err)
	page.Find("#gaia_loginform").Submit()
	a.Nil(err)

	time.Sleep(time.Second * 1)
	page.Find("#Passwd").Fill(os.Getenv("E2E_GOOGLE_PASSWORD"))
	a.Nil(err)
	page.Find("#gaia_loginform").Submit()
	a.Nil(err)

	time.Sleep(time.Second * 3)
	err = page.Find("#submit_approve_access").Click()
	a.Nil(err)
	//time.Sleep(time.Second * 10)
	// TODO: Check HTML content

	cookies, err := page.GetCookies()
	a.Nil(err)
	apiToken := getAPIToken(cookies)
	a.NotEmpty(apiToken)

	user, err := model.NewUserService(db).FindByUserAPIToken(apiToken)
	a.Nil(err)
	a.Equal(googleAccount, user.Email.Raw())
}

func TestOAuthGoogleLogout(t *testing.T) {
	if os.Getenv("CIRCLECI") != "" {
		t.Skipf("Skip because PhantomJS can't SetCookie.")
	}

	a := assert.New(t)

	_, apiToken, err := createUserAndLogin("oinume", randomEmail("oinume"), util.RandomString(16))
	a.Nil(err)

	driver := newWebDriver()
	a.Nil(driver.Start())
	defer driver.Stop()

	page, err := driver.NewPage()
	a.Nil(err)
	a.Nil(page.Navigate(server.URL))
	u, err := url.Parse(server.URL)
	a.Nil(err)
	cookie := &http.Cookie{
		Name:     controller.APITokenCookieName,
		Domain:   u.Host,
		Value:    apiToken,
		Path:     "/",
		Expires:  time.Now().Add(model.UserAPITokenExpiration),
		HttpOnly: false,
	}
	a.Nil(page.SetCookie(cookie))
	a.Nil(page.Navigate(server.URL + "/me"))

	a.Nil(page.Navigate(server.URL + "/me/logout"))
	userAPITokenService := model.NewUserAPITokenService(db)
	_, err = userAPITokenService.FindByPK(apiToken)
	a.IsType(&errors.NotFound{}, err)
}

func getAPIToken(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "apiToken" {
			return cookie.Value
		}
	}
	return ""
}
