# build stage
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git linux-headers
WORKDIR /go/src/github.com/oliviermichaelis/home-sensor
# TODO reduce go get, don't go install but go build selective files
RUN go get -t "github.com/d2r2/go-i2c" "github.com/d2r2/go-logger" "github.com/d2r2/go-bsbmp" "github.com/streadway/amqp" "github.com/influxdata/influxdb1-client/v2"
COPY . .
RUN go install -ldflags="-s -w" -v ./...    # Linker flags to reduce binary filesize. Stripping debug information

# final stage
FROM alpine

RUN addgroup -S apiserver && adduser -S apiserver -G apiserver
USER apiserver

WORKDIR /app
COPY --from=build-env /go/bin/apiserver /app/

ENTRYPOINT ./apiserver