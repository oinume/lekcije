package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/controller"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = time.UTC
var _ = fmt.Print

func TestOAuthGoogleLogin(t *testing.T) {
	if os.Getenv("CIRCLECI") != "" {
		t.Skipf("Skip because it can't render Google login page.")
	}
	a := assert.New(t)
	r := require.New(t)
	driver := newWebDriver()
	err := driver.Start()
	r.NoError(err)
	defer driver.Stop()

	page, err := driver.NewPage()
	r.NoError(err)
	r.NoError(page.Navigate(server.URL))
	//time.Sleep(15 * time.Second)

	// Check trackingId is set
	cookies, err := page.GetCookies()
	r.NoError(err)
	trackingIDCookie, err := getCookie(cookies, controller.TrackingIDCookieName)
	fmt.Printf("trackingId = %v\n", trackingIDCookie.Value)
	r.NoError(err)
	a.NotEmpty(trackingIDCookie.Value)

	signupButton := page.FindByXPath("//a[contains(@class, 'button-signup')]")
	signupURL, err := signupButton.Attribute("href")
	r.NoError(err)

	r.NoError(page.Navigate(signupURL))
	buttonGoogle := page.FindByXPath("//button[contains(@class, 'button-google')]")
	r.NoError(buttonGoogle.Click())
	//time.Sleep(15 * time.Second)

	time.Sleep(time.Second * 1)
	googleAccount := os.Getenv("E2E_GOOGLE_ACCOUNT")
	err = page.FindByXPath("//input[@name='identifier']").Fill(googleAccount)
	r.NoError(err)
	err = page.FindByXPath("//div[@id='identifierNext']/content/span").Click()
	r.NoError(err)

	time.Sleep(time.Second * 1)
	err = page.FindByXPath("//input[@name='password']").Fill(os.Getenv("E2E_GOOGLE_PASSWORD"))
	r.NoError(err)
	err = page.FindByXPath("//div[@id='passwordNext']/content/span").Click()
	r.NoError(err)

	time.Sleep(time.Second * 3)

	cookies, err = page.GetCookies()
	r.NoError(err)
	apiToken := getAPIToken(cookies)
	a.NotEmpty(apiToken)

	user, err := model.NewUserService(db).FindByUserAPIToken(apiToken)
	r.NoError(err)
	a.Equal(googleAccount, user.Email)

	// TODO: Check HTML content
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
		Domain:   strings.Split(u.Host, ":")[0], // Remove port
		Value:    apiToken,
		Path:     "/",
		Expires:  time.Now().Add(model.UserAPITokenExpiration),
		HttpOnly: false,
	}
	a.Nil(page.SetCookie(cookie))
	a.Nil(page.Navigate(server.URL + "/me"))
	//time.Sleep(time.Second * 5)

	a.Nil(page.Navigate(server.URL + "/me/logout"))
	userAPITokenService := model.NewUserAPITokenService(db)
	_, err = userAPITokenService.FindByPK(apiToken)
	a.True(errors.IsNotFound(err))
}

func getAPIToken(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "apiToken" {
			return cookie.Value
		}
	}
	return ""
}
