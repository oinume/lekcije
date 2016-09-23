package bootstrap

import (
	"reflect"
	"fmt"
	"os"
	"log"
	"time"
	"net/http"
)

type Environments struct {
	DBURL string `env:"DB_URL"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	NodeEnv string `env:"NODE_ENV"`
	RedisURL string `env:"REDIS_URL"`
}

var _ = fmt.Print
var Envs = &Environments{}

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
}

func CheckEnvs() {
	reflectValue := reflect.Indirect(reflect.ValueOf(Envs))
	for i := 0; i < reflectValue.Type().NumField(); i++ {
		fieldType := reflectValue.Type().Field(i)
		envName := fieldType.Tag.Get("env")
		envValue := os.Getenv(envName)
		if envValue == "" {
			log.Fatalf("Env '%v' must be defined.", envName)
		}
		reflectValue.FieldByName(fieldType.Name).SetString(envValue)
	}
}
