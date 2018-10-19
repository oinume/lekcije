FROM golang:1.11.0-alpine AS builder

RUN set -eux \
	&& apk --update add --no-cache \
		bash \
		make

WORKDIR /go/src/github.com/oinume/lekcije
ADD . /go/src/github.com/oinume/lekcije
RUN make all


FROM alpine
COPY --from=builder /go/src/github.com/oinume/lekcije/bin/lekcije /lekcije

ENV GCP_PROJECT_ID dummy
ENV PORT 4001
ENV GRPC_PORT 4002
ENV GOOGLE_CLIENT_ID dummy
ENV GOOGLE_CLIENT_SECRET dummy
ENV GOOGLE_ANALYTICS_ID UA-2241989-19
ENV LEKCIJE_ENV local
ENV MYSQL_USER lekcije
ENV MYSQL_PASSWORD lekcije
ENV MYSQL_HOST 192.168.99.100
ENV MYSQL_PORT 13306
ENV MYSQL_DATABASE lekcije
ENV NODE_ENV local
ENV ENCRYPTION_KEY=dummy
ENV SENDGRID_API_KEY dummy
ENV REDIS_URL redis://h:@192.168.99.100:16379
ENV E2E_GOOGLE_ACCOUNT ...
ENV E2E_GOOGLE_PASSWORD ...
ENV E2E_WEB_DRIVER chromedriver
ENV DEBUG_KEY dummy
ENV VERSION_HASH _version_
ENV ROLLBAR_ACCESS_TOKEN dummy
ENV GCP_SERVICE_ACCOUNT_KEY dummy

RUN set -eux \
	&& apk --update add --no-cache ca-certificates

EXPOSE ${HTTP_PORT:-4001}
EXPOSE ${GRPC_PORT:-4002}

ENTRYPOINT ["/lekcije"]
