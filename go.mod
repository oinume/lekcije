// +heroku goVersion go1.14
// +heroku install ./server/cmd/...
module github.com/oinume/lekcije

go 1.14

require (
	cloud.google.com/go v0.43.0
	contrib.go.opencensus.io/exporter/stackdriver v0.12.4
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/Songmu/retry v0.0.0-20170110085223-3d913ef13826
	github.com/aws/aws-sdk-go v1.21.4 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190412130859-3b1d194e553a // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/fukata/golang-stats-api-handler v1.0.0
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-sql-driver/mysql v0.0.0-20180125054745-bc14601d1bd5
	github.com/golang/protobuf v1.4.2
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/google/uuid v0.0.0-20171129191014-dec09d789f3d
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v0.0.0-20180108155640-d0c54e68681e
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/jinzhu/gorm v0.0.0-20180210142528-85774eb9dab4
	github.com/jinzhu/inflection v0.0.0-20170102125226-1c35d901db3d // indirect
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/oinume/goenum v0.0.0-20141126043735-4c1a12f41a93
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/openzipkin/zipkin-go v0.2.0
	github.com/pkg/errors v0.8.0
	github.com/rs/cors v1.3.0
	github.com/sclevine/agouti v0.0.0-20171003013254-8cf0313221cb
	github.com/sendgrid/rest v2.4.0+incompatible
	github.com/sendgrid/sendgrid-go v3.4.1+incompatible
	github.com/stretchr/testify v1.3.0
	github.com/stvp/rollbar v0.0.0-20171113052335-4a50daf855af
	go.opencensus.io v0.22.0
	go.uber.org/atomic v1.3.1 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.7.1
	goji.io/v3 v3.0.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/text v0.3.2
	google.golang.org/api v0.7.0
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884
	google.golang.org/grpc v1.29.1
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/redis.v4 v4.2.4
	gopkg.in/xmlpath.v2 v2.0.0-20150820204837-860cbeca3ebc
)
