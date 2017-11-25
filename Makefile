APP=lekcije
VENDOR_DIR=vendor
PROTO_GEN_DIR=proto-gen
GRPC_GATEWAY_REPO=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
E2E_TEST_ARGS=-v
GO_TEST_ARGS=-v
GO_TEST_PACKAGES=$(shell go list ./... | grep -v vendor | grep -v e2e)
DB_HOST=192.168.99.100
LINT_PACKAGES=$(shell go list ./... | grep -v vendor | grep -v proto | grep -v proto-gen)
VERSION_HASH_VALUE=$(shell git rev-parse HEAD | cut -c-7)
PID=$(APP).pid

all: install

.PHONY: setup
setup: install-dep install-commands

.PHONY: install-dep
install-dep:
	dep ensure -v

.PHONY: install-commands
install-commands:
	go get bitbucket.org/liamstask/goose/cmd/goose
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get honnef.co/go/tools/cmd/staticcheck
	go get honnef.co/go/tools/cmd/gosimple
	go get honnef.co/go/tools/cmd/unused

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

build:
	go build -o bin/$(APP) github.com/oinume/lekcije/server/cmd/lekcije

clean:
	${RM} bin/$(APP)

.PHONY: proto/go
proto/go:
	rm -rf $(PROTO_GEN_DIR)/go && mkdir -p $(PROTO_GEN_DIR)/go
	protoc -I/usr/local/include -I. \
  		-I$(GOPATH)/src \
  		-I$(VENDOR_DIR)/$(GRPC_GATEWAY_REPO) \
  		--go_out=plugins=grpc:$(PROTO_GEN_DIR)/go \
  		proto/echo/v1/echo.proto
	protoc -I/usr/local/include -I. -I$(GOPATH)/src -I$(VENDOR_DIR)/$(GRPC_GATEWAY_REPO) \
		--grpc-gateway_out=logtostderr=true:$(PROTO_GEN_DIR)/go \
		proto/echo/v1/echo.proto

.PHONY: ngrok
ngrok:
	ngrok http -subdomain=lekcije -host-header=localhost 4000

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
	MINIFY=true VERSION_HASH=$(VERSION_HASH_VALUE) npm run build

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
