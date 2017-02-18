E2E_TEST_ARGS=-v
GO_TEST_ARGS=-v
GO_TEST_PACKAGES=$(shell glide novendor | grep -v e2e)
DB_HOST=192.168.99.100
LINT_PACKAGES=$(shell glide novendor)

all: install

.PHONY: install

setup:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports
	glide install
	go install ./vendor/bitbucket.org/liamstask/goose/cmd/goose
	go install ./vendor/github.com/cespare/reflex
	go install ./vendor/honnef.co/go/tools/cmd/staticcheck

serve:
	go run server/cmd/lekcije/main.go

install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

e2e_test:
	go test $(E2E_TEST_ARGS) github.com/oinume/lekcije/e2e

go_test:
	go test -race $(GO_TEST_ARGS) $(GO_TEST_PACKAGES)

goimports:
	goimports -w ./server ./e2e

go_lint:
	go vet -v $(LINT_PACKAGES)
#	staticcheck $(LINT_PACKAGES)

minify_static:
	MINIFY=true VERSION_HASH=$(shell git rev-parse HEAD) npm run build

reset_db:
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije_test"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot < db/create_database.sql

