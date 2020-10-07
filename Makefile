APP = lekcije
COMMANDS = crawler daily_reporter follow_reminder notifier server teacher_error_resetter
BASE_DIR = github.com/oinume/lekcije
VENDOR_DIR = vendor
PROTO_GEN_DIR = proto-gen
GRPC_GATEWAY_REPO = github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GO_GET ?= go get
GO_TEST ?= go test -v -race -p=1 # To avoid database operations conflict
GO_TEST_E2E ?= go test -v -p=1
GO_TEST_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v e2e)
DB_HOST = 192.168.99.100
LINT_PACKAGES = $(shell go list ./... | grep -v proto-gen/go)
IMAGE_TAG ?= latest
VERSION_HASH_VALUE = $(shell git rev-parse HEAD | cut -c-7)
PID = $(APP).pid


all: build

.PHONY: setup
setup: install-tools

.PHONY: install-tools
install-tools:
	cd tools && ./install-tools.sh

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

.PHONY: build
build: $(foreach command,$(COMMANDS),build/$(command))

# TODO: find server/cmd -type d | xargs basename
# OR CLIENTS=hoge fuga proto: $(foreach var,$(CLIENTS),proto/$(var))
build/%:
	GO111MODULE=on go build -o bin/lekcije_$* $(BASE_DIR)/server/cmd/$*

clean:
	${RM} $(foreach command,$(COMMANDS),bin/lekcije_$(command))

.PHONY: db/goose/%
db/goose/%:
	goose -dir ./db/migrations mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)?charset=utf8mb4&parseTime=true&loc=UTC" $*

.PHONY: proto/go
proto/go:
	rm -rf $(PROTO_GEN_DIR)/go && mkdir -p $(PROTO_GEN_DIR)/go
	protoc -I/usr/local/include -I. \
  		-I$(GOPATH)/src \
  		-I$(VENDOR_DIR)/$(GRPC_GATEWAY_REPO) \
  		--go_out=plugins=grpc:$(PROTO_GEN_DIR)/go \
  		proto/api/v1/*.proto
	protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(VENDOR_DIR)/$(GRPC_GATEWAY_REPO) \
		--grpc-gateway_out=logtostderr=true:$(PROTO_GEN_DIR)/go \
		proto/api/v1/*.proto

.PHONY: ngrok
ngrok:
	ngrok http -subdomain=lekcije -host-header=localhost 4000

.PHONY: test
test: go-test e2e-test

.PHONY: e2e-test
e2e-test: minify-static-development
	$(GO_TEST_E2E) github.com/oinume/lekcije/e2e

.PHONY: test/db/goose/%
test/db/goose/%:
	goose -dir ./db/migrations mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)_test?charset=utf8mb4&parseTime=true&loc=UTC" $*

.PHONY: go-test
go-test:
	$(GO_TEST) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w ./server ./e2e

.PHONY: go-lint
go-lint:
	golangci-lint version
	golangci-lint run -j 4 --out-format=line-number ./...

.PHONY: docker/build/server
docker/build/server:
	docker build --pull -f docker/Dockerfile-server \
	--tag asia.gcr.io/oinume-lekcije/server:$(IMAGE_TAG) .

.PHONY: minify-static-development
minify-static-development:
	MINIFY=true VERSION_HASH=_version_ npm run build
	@echo "./static/_version_ created"

.PHONY: minify-static
minify-static:
	MINIFY=true VERSION_HASH=$(VERSION_HASH_VALUE) npm run build
	@echo "./static/$(VERSION_HASH_VALUE) created"

.PHONY: print-version-hash
print-version-hash:
	@echo $(VERSION_HASH_VALUE)

.PHONY: reset-db
reset-db:
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS lekcije_test"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot < db/create_database.sql

kill:
	kill `cat $(PID)` 2> /dev/null || true

restart: kill clean build/server
	bin/$(APP)_server & echo $$! > $(PID)

watch: restart
	fswatch -o -e ".*" -e vendor -e node_modules -e .venv -i "\\.go$$" . | xargs -n1 -I{} make restart || make kill

.PHONY: help
help:  ## show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[\/a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
