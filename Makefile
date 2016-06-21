all: dep

.PHONY: dep

dep:
	go get -u github.com/Masterminds/glide
	go get -u golang.org/x/tools/cmd/goimports
	glide install
