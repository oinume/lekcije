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
	"github.com/sclevine/agouti"
	"go.uber.org/zap/zapcore"
	"goji.io/v3"

	"github.com/oinume/lekcije/backend/config"
	"github.com/oinume/lekcije/backend/infrastructure/ga_measurement"
	interfaces "github.com/oinume/lekcije/backend/interface"
	interfaces_http "github.com/oinume/lekcije/backend/interface/http"
	"github.com/oinume/lekcije/backend/logger"
	"github.com/oinume/lekcije/backend/model"
)

var (
	server *httptest.Server
	client = http.DefaultClient
	helper = model.NewTestHelper()
	db     *gorm.DB
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	if err := os.Setenv("MYSQL_DATABASE", "lekcije_test"); err != nil {
		panic(err)
	}

	var accessLogBuffer, appLogBuffer bytes.Buffer
	appLogLevel := zapcore.InfoLevel
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		appLogLevel = logger.NewLevel(level)
	}
	args := &interfaces.ServerArgs{
		AccessLogger:        logger.NewAccessLogger(&accessLogBuffer),
		AppLogger:           logger.NewAppLogger(&appLogBuffer, appLogLevel),
		GormDB:              helper.DB(nil),
		GAMeasurementClient: ga_measurement.NewFakeClient(),
	}
	s := interfaces_http.NewServer(args, nil)
	mux := goji.NewMux()
	s.Setup(mux)
	port := config.DefaultVars.HTTPPort + 1
	server = newTestServer(mux, port)
	fmt.Printf("Test HTTP server created: port=%d, url=%s\n", port, server.URL)
	defer server.Close()

	helper.TruncateAllTables(nil)

	client.Timeout = 5 * time.Second
	if err := os.Chdir("../"); err != nil {
		panic(fmt.Errorf("os.Chdir failed: %v", err))
	}
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
