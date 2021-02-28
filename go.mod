// +heroku goVersion go1.16
// +heroku install ./backend/cmd/...
module github.com/oinume/lekcije

go 1.16

require (
	cloud.google.com/go v0.63.0
	cloud.google.com/go/storage v1.10.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.3
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/99designs/gqlgen v0.13.0 // indirect
	github.com/Songmu/retry v0.1.0
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.34.5 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec // indirect
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1
	github.com/grpc-ecosystem/grpc-gateway v1.15.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.8.0 // indirect
	github.com/matryer/moq v0.2.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oinume/goenum v0.0.0-20141126043735-4c1a12f41a93
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.3
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sclevine/agouti v3.0.0+incompatible
	github.com/sendgrid/rest v2.6.0+incompatible
	github.com/sendgrid/sendgrid-go v3.6.1+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/stvp/rollbar v0.0.0-20171113052335-4a50daf855af
	github.com/urfave/cli/v2 v2.3.0 // indirect
	github.com/vektah/dataloaden v0.3.0 // indirect
	go.opencensus.io v0.22.4
	go.uber.org/zap v1.15.0
	goji.io/v3 v3.0.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20210228012217-479acdf4ea46 // indirect
	golang.org/x/text v0.3.3
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/api v0.30.0
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70
	google.golang.org/grpc v1.32.0
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/redis.v4 v4.2.4
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
)
