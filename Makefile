all: dep

.PHONY: dep

dep:
	go get -u github.com/Masterminds/glide
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/cespare/reflex
	glide install

serve:
	go run server/main.go
