// +heroku goVersion go1.16
// +heroku install ./backend/cmd/...
module github.com/oinume/lekcije

go 1.16

require (
	cloud.google.com/go/monitoring v0.1.0 // indirect
	cloud.google.com/go/profiler v0.1.0
	cloud.google.com/go/storage v1.10.0
	cloud.google.com/go/trace v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.8
	contrib.go.opencensus.io/exporter/zipkin v0.1.2
	github.com/Songmu/retry v0.1.0
	github.com/friendsofgo/errors v0.9.2
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.4 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.8.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oinume/goenum v0.0.0-20141126043735-4c1a12f41a93
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/rollbar/rollbar-go v1.4.1
	github.com/rs/cors v1.8.0
	github.com/sclevine/agouti v3.0.0+incompatible
	github.com/sendgrid/rest v2.6.4+incompatible
	github.com/sendgrid/sendgrid-go v3.10.0+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/twitchtv/twirp v8.1.0+incompatible
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/sqlboiler/v4 v4.6.0
	github.com/volatiletech/strmangle v0.0.1
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.19.0
	goji.io/v3 v3.0.0
	golang.org/x/net v0.0.0-20210917163549-3c21e5b27794
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.7
	google.golang.org/api v0.57.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)
