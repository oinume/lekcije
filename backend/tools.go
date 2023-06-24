//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/Khan/genqlient"
	_ "github.com/matryer/moq"
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql"
)
