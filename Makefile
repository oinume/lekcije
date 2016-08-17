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
	go test -v github.com/oinume/lekcije/e2e
