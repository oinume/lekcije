all: setup

.PHONY: setup

setup:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports
	glide install
	go install ./vendor/bitbucket.org/liamstask/goose/cmd/goose
	go install ./vendor/github.com/cespare/reflex

serve:
	go run server/cmd/main.go
