FROM golang:alpine as builder

RUN     apk add --no-cache git gcc libc-dev
RUN     go get github.com/Kong/go-pluginserver

# RUN     mkdir -p /go/src/local.packages
# COPY    ./lib/ /go/src/local.packages/

RUN     mkdir /go-plugins
COPY    ./go-keycloak.go /go-plugins/go-keycloak.go
RUN     go build -buildmode plugin -o /go-plugins/go-keycloak.so /go-plugins/go-keycloak.go

FROM    kong:2.0.1-alpine

COPY    --from=builder /go/bin/go-pluginserver /usr/local/bin/go-pluginserver
RUN     mkdir /tmp/go-plugins
COPY    --from=builder /go-plugins/go-keycloak.so /tmp/go-plugins/go-keycloak.so
COPY    config.yml /tmp/config.yml

USER    root
RUN     chmod -R 777 /tmp
USER    kong