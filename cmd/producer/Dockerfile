# build stage
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git
WORKDIR /go/src/github.com/oliviermichaelis/home-sensor
RUN go get -t "github.com/d2r2/go-i2c" "github.com/d2r2/go-logger" "github.com/d2r2/go-bsbmp" "github.com/streadway/amqp" "github.com/influxdata/influxdb1-client/v2"
COPY . .
RUN go install -ldflags="-s -w" -v ./...    # Linker flags to reduce binary filesize. Stripping debug information

# final stage
FROM alpine

#RUN addgroup -S producer && adduser -S producer -G producer
#USER producer

WORKDIR /app
COPY --from=build-env /go/bin/producer /app/

ENTRYPOINT ./producer