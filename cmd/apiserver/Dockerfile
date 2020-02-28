# build stage
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git linux-headers
WORKDIR /go/src/github.com/oliviermichaelis/home-sensor
#RUN go get -t "github.com/influxdata/influxdb1-client/v2"
COPY . .
RUN go install -ldflags="-s -w" -v ./...    # Linker flags to reduce binary filesize. Stripping debug information

# final stage
FROM alpine

RUN addgroup -S apiserver && adduser -S apiserver -G apiserver
USER apiserver

WORKDIR /app
COPY --from=build-env /go/bin/apiserver /app/

ENTRYPOINT ./apiserver
