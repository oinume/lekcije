package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	t := `
add-panic-variants = false
no-tests = true
output  = "./model2"
pkgname = "model2"

[mysql]
  dbname  = "%s"
  host    = "%s"
  port    = %s
  user    = "%s"
  pass    = "%s"
  sslmode = "false"
  blacklist = ["goose_db_version"]
`
	s := fmt.Sprintf(
		t,
		getenvDefault("MYSQL_DATABASE", "hoge"),
		getenvDefault("MYSQL_HOST", "localhost"),
		getenvDefault("MYSQL_PORT", "3306"),
		getenvDefault("MYSQL_USER", "root"),
		getenvDefault("MYSQL_PASSWORD", "root"),
	)
	s = strings.TrimSpace(s)
	fmt.Println(s)
}

func getenvDefault(name, d string) string {
	v := os.Getenv(name)
	if v == "" {
		return d
	}
	return v
}
