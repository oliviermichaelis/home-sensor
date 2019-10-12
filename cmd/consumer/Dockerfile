# build stage
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git
ADD . /src
RUN go get -t "github.com/influxdata/influxdb1-client/v2" "github.com/d2r2/go-bsbmp" "github.com/d2r2/go-logger" "github.com/streadway/amqp"
RUN cd /src && go build -o producer

# final stage
FROM alpine

# Needed to be removed due to the permissions when the device is handed into the container
#RUN addgroup -S producer && adduser -S producer -G producer
#USER producer
WORKDIR /app
COPY --from=build-env /src/producer /app/
ENTRYPOINT ./producer