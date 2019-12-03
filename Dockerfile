FROM golang:1.13.4-stretch as builder
WORKDIR /go/src/github.com/asdfsx/k8s-device-plugin/
ENV GO111MODULE=on
COPY . .
RUN make build

FROM debian:stretch-slim

ADD /bin/k8s-device-plugin /bin/k8s-device-plugin

CMD ["/bin/k8s-device-plugin"]
