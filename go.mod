// +heroku goVersion go1.16
// +heroku install ./backend/cmd/...
module github.com/oinume/lekcije

go 1.16

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/twitchtv/twirp v8.1.2+incompatible
	google.golang.org/protobuf v1.28.1
)

replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
