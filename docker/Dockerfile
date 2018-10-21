FROM golang:1.11.0-alpine AS builder

RUN set -eux \
	&& apk --update add --no-cache \
		bash \
		make

WORKDIR /go/src/github.com/oinume/lekcije
ADD . /go/src/github.com/oinume/lekcije
RUN make build/lekcije


FROM alpine
COPY --from=builder /go/src/github.com/oinume/lekcije/bin/lekcije /lekcije

RUN set -eux \
	&& apk --update add --no-cache ca-certificates

EXPOSE ${PORT:-4001} # HTTP
EXPOSE ${GRPC_PORT:-4002}

ENTRYPOINT ["/lekcije"]
