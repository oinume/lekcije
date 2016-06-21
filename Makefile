PID      = /tmp/lekcije.pid
GO_FILES = $(wildcard *.go)

all: dep

#.PHONY: dep
.PHONY: dep serve restart kill stuff # let's go to reserve rules names

dep:
	go get -u github.com/Masterminds/glide
	go get -u golang.org/x/tools/cmd/goimports
	glide install

#serve:
#	@make start
#	@fswatch -o . | xargs -n1 -I{} make restart || make kill
#
#kill:
##	@kill `cat $(PID)` || true
##	@pkill -f 'go run' || true
##	@pgrep -f go-build > hoge
#	@kill `lsof -t -i :5000` || true
#
#stuff:
#	@echo "actually do nothing"
#
#start:
#	@go run server/main.go & echo $$! > $(PID)
#
#restart:
#	@make kill
#	@make stuff
#	@go run server/main.go & echo $$! > $(PID)
