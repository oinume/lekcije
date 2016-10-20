package config

import (
	"os"
	"strconv"
	"time"
)

var (
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
)

func StaticURL() string {
	if IsProductionEnv() {
		return "https://asset.lekcije.com/static"
	} else if IsDevelopmentEnv() {
		return "http://asset.local.lekcije.com/static"
	} else {
		return "/static"
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

func LocalTimezone() *time.Location {
	return jst
}
