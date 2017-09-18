package e2e

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/oinume/lekcije/server/route"
	"github.com/sclevine/agouti"
)

var server *httptest.Server
var client = http.DefaultClient
var db *gorm.DB

func TestMain(m *testing.M) {
	bootstrap.CheckCLIEnvVars()
	dbURL := model.ReplaceToTestDBURL(bootstrap.CLIEnvVars.DBURL())
	if err := os.Setenv("MYSQL_DATABASE", "lekcije_test"); err != nil {
		// TODO: Not use panic
		panic(err)
	}
	bootstrap.CheckServerEnvVars()

	var accessLogBuffer, appLogBuffer bytes.Buffer
	logger.InitializeAccessLogger(&accessLogBuffer)
	logger.InitializeAppLogger(&appLogBuffer)

	var err error
	db, err = model.OpenDB(dbURL, 1, true) // TODO: env
	if err != nil {
		panic(err)
	}
	if err := model.TruncateAllTables(db, model.GetDBName(dbURL)); err != nil {
		panic(err)
	}

	port := config.ListenPort()
	routes := route.Create(nil) // TODO: grpc-gateway
	port += 1
	server = newTestServer(routes, port)
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
	e2eWebDriver := os.Getenv("E2E_WEB_DRIVER")
	var driver *agouti.WebDriver
	switch strings.ToLower(e2eWebDriver) {
	case "chromedriver":
		driver = agouti.ChromeDriver()
	case "phantomjs":
		driver = agouti.PhantomJS()
	default:
		driver = agouti.PhantomJS()
	}
	driver.HTTPClient = client
	return driver
}
