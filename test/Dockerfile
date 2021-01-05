FROM golang:alpine

RUN     apk add --no-cache gcc musl-dev
COPY    ./functionalTest/ /go/src

WORKDIR /go/src

ENTRYPOINT [ "go", "test" ]