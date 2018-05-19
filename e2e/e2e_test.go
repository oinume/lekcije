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
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/interfaces"
	interfaces_http "github.com/oinume/lekcije/server/interfaces/http"
	"github.com/oinume/lekcije/server/logger"
	"github.com/oinume/lekcije/server/model"
	"github.com/sclevine/agouti"
	"go.uber.org/zap/zapcore"
)

var server *httptest.Server
var client = http.DefaultClient
var db *gorm.DB
var helper = model.NewTestHelper()

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	db = helper.DB()
	if err := os.Setenv("MYSQL_DATABASE", "lekcije_test"); err != nil {
		// TODO: Not use panic
		panic(err)
	}

	var accessLogBuffer, appLogBuffer bytes.Buffer
	logger.InitializeAccessLogger(&accessLogBuffer)
	appLogLevel := zapcore.InfoLevel
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		appLogLevel = logger.NewLevel(level)
	}
	logger.InitializeAppLogger(&appLogBuffer, appLogLevel)

	helper.TruncateAllTables(helper.DB())

	port := config.DefaultVars.HTTPPort
	args := &interfaces.ServerArgs{
		DB: db,
		//RedisClient: redis
	}
	s := interfaces_http.NewServer(args)
	routes := s.CreateRoutes(nil) // TODO: grpc-gateway
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
	case "chromedriver_headless":
		driver = agouti.ChromeDriver(
			agouti.ChromeOptions("args", []string{
				"--headless",             // headlessモードの指定
				"--window-size=1280,800", // ウィンドウサイズの指定
			}),
			agouti.Debug,
		)
	default:
		driver = agouti.ChromeDriver()
	}
	driver.HTTPClient = client
	return driver
}
