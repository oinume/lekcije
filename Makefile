E2E_TEST_ARGS=-v
GO_TEST_ARGS=-v
GO_TEST_PACKAGES=$(shell glide novendor | grep -v e2e)
DB_HOST=192.168.99.100
LINT_PACKAGES=$(shell glide novendor)
VERSION_HASH=$(shell git rev-parse HEAD | cut -c-7)

all: install

.PHONY: setup
setup:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports
	glide install
	go install ./vendor/bitbucket.org/liamstask/goose/cmd/goose
	go install ./vendor/github.com/cespare/reflex
	go install ./vendor/honnef.co/go/tools/cmd/staticcheck
	go install ./vendor/honnef.co/go/tools/cmd/gosimple
	go install ./vendor/honnef.co/go/tools/cmd/unused

.PHONY: serve
serve:
	go run server/cmd/lekcije/main.go

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

.PHONY: test
test: go_test e2e_test

.PHONY: e2e_test
e2e_test: minify_static_development
	go test $(E2E_TEST_ARGS) github.com/oinume/lekcije/e2e

.PHONY: go_test
go_test:
	go test $(GO_TEST_ARGS) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w ./server ./e2e

.PHONY: go_lint
go_lint: go_vet go_staticcheck go_simple

.PHONY: go_vet
go_vet:
	go vet -v $(LINT_PACKAGES)

.PHONY: go_staticcheck
go_staticcheck:
	staticcheck $(LINT_PACKAGES)

.PHONY: go_simple
go_simple:
	gosimple $(LINT_PACKAGES)

.PHONY: minify_static_development
minify_static_development:
	MINIFY=true VERSION_HASH=_version_ npm run build

.PHONY: minify_static
minify_static:
	MINIFY=true VERSION_HASH=$(VERSION_HASH) npm run build

.PHONY: print_version_hash
print_version_hash:
	@echo $(VERSION_HASH)

.PHONY: reset_db
reset_db:
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije_test"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot < db/create_database.sql
