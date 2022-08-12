// +heroku goVersion go1.16
// +heroku install ./backend/cmd/...
module github.com/oinume/lekcije

go 1.16

require (
	cloud.google.com/go/profiler v0.3.0
	cloud.google.com/go/storage v1.24.0
	github.com/99designs/gqlgen v0.17.13
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.8.4
	github.com/Songmu/retry v0.1.0
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/ericlagergren/decimal v0.0.0-20211103172832-aca2edc11f73
	github.com/friendsofgo/errors v0.9.2
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.8
	github.com/google/uuid v1.3.0
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/kr/text v0.2.0 // indirect
	github.com/matryer/moq v0.2.7
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/oinume/goenum v0.0.0-20141126043735-4c1a12f41a93
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.5.3
	github.com/rollbar/rollbar-go v1.4.4
	github.com/rs/cors v1.8.2
	github.com/sendgrid/rest v2.6.5+incompatible
	github.com/sendgrid/sendgrid-go v3.10.3+incompatible
	github.com/sethvargo/go-envconfig v0.8.1
	github.com/stretchr/testify v1.8.0
	github.com/twitchtv/twirp v8.1.2+incompatible
	github.com/urfave/cli v1.22.2 // indirect
	github.com/vektah/gqlparser/v2 v2.4.7
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/sqlboiler/v4 v4.12.0
	github.com/volatiletech/strmangle v0.0.4
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.33.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.33.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/jaeger v1.8.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.8.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/zap v1.22.0
	goji.io/v3 v3.0.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.11 // indirect
	google.golang.org/api v0.91.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)

replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
