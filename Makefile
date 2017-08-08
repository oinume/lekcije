E2E_TEST_ARGS=-v
GO_TEST_ARGS=-v
GO_TEST_PACKAGES=$(shell glide novendor | grep -v e2e)
DB_HOST=192.168.99.100
LINT_PACKAGES=$(shell glide novendor)
VERSION_HASH=$(shell git rev-parse HEAD | cut -c-7)

all: install

.PHONY: setup
setup: install-glide install-dep install-commands

.PHONY: install-glide
install-glide:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports

.PHONY: install-dep
install-dep:
	glide install

.PHONY: install-commands
install-commands:
	go install ./vendor/bitbucket.org/liamstask/goose/cmd/goose
	go install ./vendor/github.com/cespare/reflex
	go install ./vendor/honnef.co/go/tools/cmd/staticcheck
	go install ./vendor/honnef.co/go/tools/cmd/gosimple
	go install ./vendor/honnef.co/go/tools/cmd/unused

.PHONY: serve
serve:
	go run server/cmd/lekcije/main.go

.PHONY: reflex
reflex:
	reflex -R node_modules -R vendor -R .venv -r '\.go$$' -s make serve

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

.PHONY: test
test: go-test e2e-test

.PHONY: e2e-test
e2e-test: minify-static-development
	go test $(E2E_TEST_ARGS) github.com/oinume/lekcije/e2e

.PHONY: go-test
go-test:
	go test $(GO_TEST_ARGS) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w ./server ./e2e

.PHONY: go-lint
go-lint: go-vet go-staticcheck go-simple

.PHONY: go-vet
go-vet:
	go vet -v $(LINT_PACKAGES)

.PHONY: go-staticcheck
go-staticcheck:
	staticcheck $(LINT_PACKAGES)

.PHONY: go-simple
go-simple:
	gosimple $(LINT_PACKAGES)

.PHONY: minify-static-development
minify-static-development:
	MINIFY=true VERSION_HASH=_version_ npm run build

.PHONY: minify-static
minify-static:
	MINIFY=true VERSION_HASH=$(shell git rev-parse HEAD | cut -c-7) npm run build

.PHONY: print-version-hash
print-version-hash:
	@echo $(VERSION_HASH)

.PHONY: reset-db
reset-db:
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije_test"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot < db/create_database.sql
