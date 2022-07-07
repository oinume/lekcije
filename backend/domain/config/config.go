package config

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	envconfig "github.com/sethvargo/go-envconfig"
)

const (
	DefaultTracerName = "github.com/oinume/lekcije"
)

var (
	asiaTokyo = time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp = time.Now().UTC()
)

type MySQL struct {
	User     string `env:"MYSQL_USER"`
	Password string `env:"MYSQL_PASSWORD"`
	Host     string `env:"MYSQL_HOST"`
	Port     string `env:"MYSQL_PORT"`
	Database string `env:"MYSQL_DATABASE"`
}

type Trace struct {
	Enable       bool    `env:"TRACE_ENABLED"`
	Exporter     string  `env:"TRACE_EXPORTER"`
	SamplingRate float64 `env:"TRACE_SAMPLING_RATE"`
	ExporterURL  string  `env:"TRACE_EXPORTER_URL"`
}

type Vars struct {
	*MySQL
	*Trace
	EncryptionKey             string `env:"ENCRYPTION_KEY"`
	NodeEnv                   string `env:"NODE_ENV"`
	ServiceEnv                string `env:"LEKCIJE_ENV" required:"true"`
	GCPProjectID              string `env:"GCP_PROJECT_ID"`
	GCPServiceAccountKey      string `env:"GCP_SERVICE_ACCOUNT_KEY"`
	EnableFetcherHTTP2        bool   `env:"ENABLE_FETCHER_HTTP2" default:"true"`
	EnableStackdriverProfiler bool   `env:"ENABLE_STACKDRIVER_PROFILER"`
	GoogleClientID            string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret        string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleAnalyticsID         string `env:"GOOGLE_ANALYTICS_ID"`
	HTTPPort                  int    `env:"PORT" default:"4001"`
	RollbarAccessToken        string `env:"ROLLBAR_ACCESS_TOKEN"`
	VersionHash               string `env:"VERSION_HASH"`
	DebugSQL                  bool   `env:"DEBUG_SQL"`
	LocalLocation             *time.Location
}

func Process() (*Vars, error) {
	var vars Vars
	if err := envconfig.Process(context.Background(), &vars); err != nil {
		return nil, err
	}
	vars.LocalLocation = asiaTokyo
	if vars.VersionHash == "" {
		vars.VersionHash = timestamp.Format("20060102150405")
	}
	return &vars, nil
}

func MustProcess() *Vars {
	vars, err := Process()
	if err != nil {
		panic(err)
	}
	return vars
}

var DefaultVars = &Vars{}
var once sync.Once

func MustProcessDefault() {
	once.Do(func() {
		DefaultVars = MustProcess()
	})
}

func (v *Vars) DBURL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		v.MySQL.User, v.MySQL.Password, v.MySQL.Host, v.MySQL.Port, v.MySQL.Database,
	)
}

func (v *Vars) StaticURL() string {
	if v.IsProductionEnv() {
		return "https://asset.lekcije.com/static/" + v.VersionHash
	} else if v.IsDevelopmentEnv() {
		return "http://asset.local.lekcije.com/static/" + v.VersionHash
	} else {
		return "/static/" + v.VersionHash
	}
}

func (v *Vars) WebURL() string {
	if v.IsProductionEnv() {
		return "https://www.lekcije.com"
	} else if v.IsDevelopmentEnv() {
		return "http://www.local.lekcije.com"
	} else {
		return "http://localhost:4000"
	}
}

func (v *Vars) IsProductionEnv() bool {
	return v.ServiceEnv == "production"
}

func (v *Vars) IsDevelopmentEnv() bool {
	return v.ServiceEnv == "development"
}

func (v *Vars) IsLocalEnv() bool {
	return v.ServiceEnv == "local"
}

func (v *Vars) WebURLScheme(r *http.Request) string {
	if v.IsProductionEnv() {
		return "https"
	}
	if r != nil && r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

func StaticURL() string {
	return DefaultVars.StaticURL()
}

func WebURL() string {
	return DefaultVars.WebURL()
}

func WebURLScheme(r *http.Request) string {
	return DefaultVars.WebURLScheme(r)
}

func LocalLocation() *time.Location {
	return DefaultVars.LocalLocation
}

func IsDevelopmentEnv() bool {
	return DefaultVars.IsDevelopmentEnv()
}

func IsLocalEnv() bool {
	return DefaultVars.IsLocalEnv()
}

func IsProductionEnv() bool {
	return DefaultVars.IsProductionEnv()
}
