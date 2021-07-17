package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oinume/lekcije/backend/errors"
	interfaces_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/model"
	"github.com/oinume/lekcije/backend/util"
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
	defer func() { _ = driver.Stop() }()

	page, err := driver.NewPage()
	r.NoError(err)
	r.NoError(page.Navigate(server.URL))
	//time.Sleep(15 * time.Second)

	// Check trackingId is set
	cookies, err := page.GetCookies()
	r.NoError(err)
	trackingIDCookie, err := getCookie(cookies, interfaces_http.TrackingIDCookieName)
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

	time.Sleep(time.Second * 3)
	err = page.FindByXPath("//input[@name='password']").Fill(os.Getenv("E2E_GOOGLE_PASSWORD"))
	r.NoError(err)
	err = page.FindByXPath("//div[@id='passwordNext']/content/span").Click()
	r.NoError(err)

	time.Sleep(time.Second * 4)

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
	r := require.New(t)

	_, apiToken, err := createUserAndLogin("oinume", randomEmail("oinume"), util.RandomString(16))
	r.NoError(err)

	driver := newWebDriver()
	r.NoError(driver.Start())
	defer func() { _ = driver.Stop() }()

	page, err := driver.NewPage()
	r.NoError(err)
	r.NoError(page.Navigate(server.URL))
	u, err := url.Parse(server.URL)
	r.NoError(err)
	cookie := &http.Cookie{
		Name:     interfaces_http.APITokenCookieName,
		Domain:   strings.Split(u.Host, ":")[0], // Remove port
		Value:    apiToken,
		Path:     "/",
		Expires:  time.Now().Add(model.UserAPITokenExpiration),
		HttpOnly: false,
	}
	r.NoError(page.SetCookie(cookie))
	r.NoError(page.Navigate(server.URL + "/me"))
	time.Sleep(2 * time.Second)

	r.NoError(page.Navigate(server.URL + "/me/logout"))
	time.Sleep(2 * time.Second)

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
