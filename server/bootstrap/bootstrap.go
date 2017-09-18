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

//$MYSQLDUMP -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST $MYSQL_DATABASE | bzip2 -9 > lekcije_$DATE.dump.bz2

type CLIEnvVarsType struct {
	MySQLUser     string `env:"MYSQL_USER"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLHost     string `env:"MYSQL_HOST"`
	MySQLPort     string `env:"MYSQL_PORT"`
	MySQLDatabase string `env:"MYSQL_DATABASE"`
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	NodeEnv       string `env:"NODE_ENV"`
	RedisURL      string `env:"REDIS_URL"`
}

func (t CLIEnvVarsType) DBURL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		t.MySQLUser, t.MySQLPassword, t.MySQLHost, t.MySQLPort, t.MySQLDatabase,
	)
}

type HTTPServerEnvVarsType struct {
	MySQLUser          string `env:"MYSQL_USER"`
	MySQLPassword      string `env:"MYSQL_PASSWORD"`
	MySQLHost          string `env:"MYSQL_HOST"`
	MySQLPort          string `env:"MYSQL_PORT"`
	MySQLDatabase      string `env:"MYSQL_DATABASE"`
	EncryptionKey      string `env:"ENCRYPTION_KEY"`
	NodeEnv            string `env:"NODE_ENV"`
	RedisURL           string `env:"REDIS_URL"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GRPCPort           string `env:"GRPC_PORT"`
}

func (t HTTPServerEnvVarsType) DBURL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		t.MySQLUser, t.MySQLPassword, t.MySQLHost, t.MySQLPort, t.MySQLDatabase,
	)
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
