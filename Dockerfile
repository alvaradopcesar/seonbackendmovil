FROM golang:1.9
# FROM imb_base:latest

# MAINTAINER Cesar Alvarado <alvaradopcesar@gmail.com>
# USER root

ENV ENVIRONMENT=xx

RUN mkdir /go/src/ProductService2

WORKDIR /go/src/ProductService2

ADD . /go/src/ProductService2

RUN go get -d -v . \
    && go build *.go

CMD ["./setup.sh"]