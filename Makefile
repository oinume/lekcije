all: setup

.PHONY: setup

setup:
	go get -u github.com/Masterminds/glide
	go get -u golang.org/x/tools/cmd/goimports
	glide install
