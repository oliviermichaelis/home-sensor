# home-sensor

[![CircleCI](https://circleci.com/gh/oliviermichaelis/home-sensor.svg?style=svg)](https://circleci.com/gh/oliviermichaelis/home-sensor)

## Setup
### Create secrets:
`kubectl create secret generic influxdata --from-literal=token='' --from-literal=org='' --from-literal=url=''`

`helm` contains the helm chart used to deploy to Kubernetes.
