# build stage
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git
ADD . /src
RUN go get -t "github.com/d2r2/go-i2c" "github.com/streadway/amqp"
RUN cd /src && go build -o consumer

# final stage
FROM alpine

RUN addgroup -S producer && adduser -S producer -G producer
USER producer
WORKDIR /app
COPY --from=build-env /src/consumer /app/
ENTRYPOINT ./consumer