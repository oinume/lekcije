APP = lekcije
BASE_DIR = github.com/oinume/lekcije
VENDOR_DIR = vendor
PROTO_GEN_DIR = proto-gen
GRPC_GATEWAY_REPO = github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GO_GET ?= go get
GO_TEST ?= go test -v -race -p=1 # To avoid database operations conflict
GO_TEST_E2E ?= go test -v -p=1
GO_TEST_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v e2e)
DB_HOST = 192.168.99.100
LINT_PACKAGES = $(shell go list ./...)
VERSION_HASH_VALUE = $(shell git rev-parse HEAD | cut -c-7)
PID = $(APP).pid


all: build

.PHONY: setup
setup: install-dep install-commands

.PHONY: install-dep
install-dep:
	dep ensure -v

.PHONY: install-commands
install-commands:
	$(GO_GET) bitbucket.org/liamstask/goose/cmd/goose
	$(GO_GET) github.com/golang/protobuf/protoc-gen-go
	$(GO_GET) github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	$(GO_GET) honnef.co/go/tools/cmd/staticcheck
	$(GO_GET) honnef.co/go/tools/cmd/gosimple
	$(GO_GET) honnef.co/go/tools/cmd/unused

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

.PHONY: build
build: build/$(APP) build/notifier

build/$(APP):
	go build -o bin/$(APP) $(BASE_DIR)/server/cmd/lekcije

build/notifier:
	go build -o bin/notifier $(BASE_DIR)/server/cmd/notifier

# TODO: build/xxx

clean:
	${RM} bin/$(APP) bin/notifier

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

.PHONY: go-test
go-test:
	$(GO_TEST) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w ./server ./e2e

.PHONY: go-lint
go-lint: go-staticcheck go-simple

.PHONY: go-vet
go-vet:
	go vet -v $(LINT_PACKAGES)

.PHONY: go-staticcheck
go-staticcheck:
	staticcheck -ignore "github.com/oinume/lekcije/proto-gen/go/proto/api/v1/*.go:SA1019" $(LINT_PACKAGES)

.PHONY: go-simple
go-simple:
	gosimple $(LINT_PACKAGES)

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

restart: kill clean build
	bin/$(APP) & echo $$! > $(PID)

watch: restart
	fswatch -o -e ".*" -e vendor -e node_modules -e .venv -i "\\.go$$" . | xargs -n1 -I{} make restart || make kill
