ARG tag=0.1

FROM golang:1.10.1 AS build
WORKDIR /go/src/github.com/vanveele/nomad-service-alerter/
COPY .
RUN go get \
    go build


FROM alpine:latest
MAINTAINER Robert van Veelen <rvanveel@citadel.com>
ENTRYPOINT ["/init"]
CMD ["/nomad-service-alerter"]

EXPOSE 8000
CMD ["/nomad-service-alerter"]
ENV nomad_server=nomad:8500\
    env=qa \
    region=primary \
    alert_switch=on \
    consul_server=localhost:8500 \
    consul_datacenter=primary \
    kafka_brokers=kafka:9092 \
    kafka_topic=nomad-service-alerts    

COPY --from=build /nomad-service-alerter
