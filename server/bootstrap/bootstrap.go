package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/stvp/rollbar"
)

// TODO: Fix reflection problem and use this struct
// http://stackoverflow.com/questions/24333494/golang-reflection-on-embedded-structs
//type CommonEnvVarsType struct {
//	DBURL string `env:"DB_URL"`
//	EncryptionKey string `env:"ENCRYPTION_KEY"`
//	NodeEnv string `env:"NODE_ENV"`
//	RedisURL string `env:"REDIS_URL"`
//}

type CLIEnvVarsType struct {
	DBURL         string `env:"DB_URL"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	NodeEnv       string `env:"NODE_ENV"`
	RedisURL      string `env:"REDIS_URL"`
}

type HTTPServerEnvVarsType struct {
	DBURL              string `env:"DB_URL"`
	EncryptionKey      string `env:"ENCRYPTION_KEY"`
	NodeEnv            string `env:"NODE_ENV"`
	RedisURL           string `env:"REDIS_URL"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

var _ = fmt.Print
var CLIEnvVars = CLIEnvVarsType{}
var ServerEnvVars = HTTPServerEnvVarsType{}

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
}

func CheckCLIEnvVars() {
	checkEnvVars(reflect.Indirect(reflect.ValueOf(&CLIEnvVars)))
	rollbar.Token = os.Getenv("ROLLBAR_ACCESS_TOKEN")
	rollbar.Endpoint = os.Getenv("ROLLBAR_ENDPOINT")
	rollbar.Environment = ServerEnvVars.NodeEnv
}

func CheckServerEnvVars() {
	checkEnvVars(reflect.Indirect(reflect.ValueOf(&ServerEnvVars)))
	rollbar.Token = os.Getenv("ROLLBAR_ACCESS_TOKEN")
	rollbar.Endpoint = os.Getenv("ROLLBAR_ENDPOINT")
	rollbar.Environment = ServerEnvVars.NodeEnv
}

func checkEnvVars(value reflect.Value) {
	for i := 0; i < value.Type().NumField(); i++ {
		fieldType := value.Type().Field(i)
		envName := fieldType.Tag.Get("env")
		if envName == "" {
			continue
		}
		envValue := os.Getenv(envName)
		if envValue == "" {
			log.Fatalf("Env '%v' must be defined.", envName)
		}
		value.FieldByName(fieldType.Name).SetString(envValue)
	}
}
