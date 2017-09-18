package config

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	jst         = time.FixedZone("Asia/Tokyo", 9*60*60)
	versionHash = os.Getenv("VERSION_HASH")
	timestamp   = time.Now().UTC()
)

func StaticURL() string {
	if IsProductionEnv() {
		return "https://asset.lekcije.com/static/" + VersionHash()
	} else if IsDevelopmentEnv() {
		return "http://asset.local.lekcije.com/static/" + VersionHash()
	} else {
		return "/static/" + VersionHash()
	}
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

func GoogleAnalyticsID() string {
	return os.Getenv("GOOGLE_ANALYTICS_ID")
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

func ListenPort() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return -1
	}
	return int(p)
}

func GRPCListenPort() int {
	port := os.Getenv("GPRC_PORT")
	if port == "" {
		port = "4002"
	}
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return -1
	}
	return int(p)
}

func LocalTimezone() *time.Location {
	return jst
}

func VersionHash() string {
	if versionHash == "" {
		return timestamp.Format("20060102150405")
	} else {
		return versionHash
	}
}

func SetVersionHash(version string) {
	versionHash = version
}
