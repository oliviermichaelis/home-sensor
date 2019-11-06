# home-sensor

[![CircleCI](https://circleci.com/gh/oliviermichaelis/home-sensor.svg?style=svg)](https://circleci.com/gh/oliviermichaelis/home-sensor)

`cmd/producer` contains the source file for the producer. The producer reads the sensor periodically and sends a message containing the measurements to a RabbitMQ queue. Both producer and consumer use a connection pool implemented in `pkg/connect/connect.go`, which reuses connections or redials in case an existing connectionw as invalidated

`cmd/consumer` contains the source file for the consumer. The consumer dequeues messages from RabbitMQ and inserts them into an InfluxDB timeseries database.

`helm` contains the helm chart used to deploy the system to Kubernetes.

`test` contains (few) unit and integration tests.
