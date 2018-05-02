package config

import (
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stvp/rollbar"
)

var (
	jst         = time.FixedZone("Asia/Tokyo", 9*60*60)
	versionHash = os.Getenv("VERSION_HASH")
	timestamp   = time.Now().UTC()
)

type Vars struct {
	MySQLUser          string `envconfig:"MYSQL_USER"`
	MySQLPassword      string `envconfig:"MYSQL_PASSWORD"`
	MySQLHost          string `envconfig:"MYSQL_HOST"`
	MySQLPort          string `envconfig:"MYSQL_PORT"`
	MySQLDatabase      string `envconfig:"MYSQL_DATABASE"`
	EncryptionKey      string `envconfig:"ENCRYPTION_KEY"`
	NodeEnv            string `envconfig:"NODE_ENV"`
	ServiceEnv         string `envconfig:"LEKCIJE_ENV" required:"true"`
	RedisURL           string `envconfig:"REDIS_URL"`
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleAnalyticsID  string `envconfig:"GOOGLE_ANALYTICS_ID"`
	HTTPPort           int    `envconfig:"PORT" default:"4001"`
	GRPCPort           int    `envconfig:"GRPC_PORT" default:"4002"`
	RollbarAccessToken string `envconfig:"ROLLBAR_ACCESS_TOKEN"`
	VersionHash        string `envconfig:"VERSION_HASH"`
	LocalTimeZone      *time.Location
}

func Process() (*Vars, error) {
	var vars Vars
	if err := envconfig.Process("", &vars); err != nil {
		return nil, err
	}

	vars.LocalTimeZone = jst
	if vars.VersionHash == "" {
		vars.VersionHash = timestamp.Format("20060102150405")
	}
	// TODO: Make it optional
	rollbar.Token = vars.RollbarAccessToken
	rollbar.Endpoint = "https://api.rollbar.com/api/1/item/"
	rollbar.Environment = vars.NodeEnv // TODO: lekcije_env

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

func MustProcessDefault() {
	DefaultVars = MustProcess()
}

func (v *Vars) StaticURL() string {
	if IsProductionEnv() {
		return "https://asset.lekcije.com/static/" + v.VersionHash
	} else if IsDevelopmentEnv() {
		return "http://asset.local.lekcije.com/static/" + v.VersionHash
	} else {
		return "/static/" + v.VersionHash
	}
}

func (v *Vars) WebURL() string {
	if IsProductionEnv() {
		return "https://www.lekcije.com"
	} else if IsDevelopmentEnv() {
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

func WebURL() string {
	if IsProductionEnv() {
		return "https://www.lekcije.com"
	} else if IsDevelopmentEnv() {
		return "http://www.local.lekcije.com"
	} else {
		return "http://localhost:4000"
	}
}

func EnvString() string {
	return os.Getenv("NODE_ENV")
}

func IsProductionEnv() bool {
	return EnvString() == "production"
}

func IsDevelopmentEnv() bool {
	return EnvString() == "development"
}

func IsLocalEnv() bool {
	return EnvString() == "local"
}

func WebURLScheme(r *http.Request) string {
	if IsProductionEnv() {
		return "https"
	}
	if r != nil && r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

func LocalTimezone() *time.Location {
	return jst
}
