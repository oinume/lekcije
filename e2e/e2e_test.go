package e2e

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/mux"
	"github.com/sclevine/agouti"
)

var server *httptest.Server
var client = http.DefaultClient

func TestMain(m *testing.M) {
	dbURL := model.ReplaceToTestDBURL(os.Getenv("DB_URL"))
	if err := os.Setenv("DB_URL", dbURL); err != nil {
		// TODO: Not use panic
		panic(err)
	}
	bootstrap.CheckHTTPServerEnvVars()

	var accessLogBuffer, appLogBuffer bytes.Buffer
	logger.InitializeAccessLogger(&accessLogBuffer)
	logger.InitializeAppLogger(&appLogBuffer)

	db, err := model.OpenDB(dbURL)
	if err != nil {
		panic(err)
	}
	if err := model.TruncateAllTables(db, model.GetDBName(dbURL)); err != nil {
		panic(err)
	}

	port := config.ListenPort()
	mux := mux.Create()
	port += 1
	server = newTestServer(mux, port)
	fmt.Printf("Test HTTP server created: port=%d, url=%s\n", port, server.URL)
	defer server.Close()

	client.Timeout = 5 * time.Second
	os.Chdir("../")
	status := m.Run()
	defer os.Exit(status)
}

// newTestServer returns a new test Server with fixed port number.
func newTestServer(handler http.Handler, port int) *httptest.Server {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		if listener, err = net.Listen("tcp6", fmt.Sprintf("[::1]:%d", port)); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	ts := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: handler},
	}
	ts.Start()
	return ts
}

func newWebDriver() *agouti.WebDriver {
	driver := agouti.ChromeDriver()
	//driver := agouti.PhantomJS()
	driver.HTTPClient = client
	return driver
}
