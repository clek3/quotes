FROM golang:latest

WORKDIR /go/src/github.com/clek3/quotes

ADD ./src /go/src/github.com/clek3/quotes/src

RUN go get ./...
