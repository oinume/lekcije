package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/matryer/moq"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql"
)
