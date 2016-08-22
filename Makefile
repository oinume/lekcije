E2E_TEST_ARGS=-v
GO_TEST_ARGS=-v
GO_TEST_PACKAGES=$(shell go list ./... | grep -v cmd | grep -v e2e | grep -v vendor)
#DB_DSN=lekcije:lekcije@tcp(192.168.99.100:13306)/lekcije

all: install

.PHONY: install

setup:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports
	glide install
	go install ./vendor/bitbucket.org/liamstask/goose/cmd/goose
	go install ./vendor/github.com/cespare/reflex

serve:
	go run server/cmd/lekcije/main.go

install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

e2e_test:
	go test $(E2E_TEST_ARGS) github.com/oinume/lekcije/e2e

go_test:
	go test $(GO_TEST_ARGS) $(GO_TEST_PACKAGES)
