FROM golang:1.13.4-stretch as builder
WORKDIR /go/src/github.com/chendotjs/k8s-device-plugin/
ENV GO111MODULE=on
COPY . .
RUN make build

FROM debian:stretch-slim
WORKDIR /bin
COPY --from=builder /go/src/github.com/chendotjs/k8s-device-plugin/bin/k8s-device-plugin .
ENTRYPOINT ["/bin/k8s-device-plugin"]
CMD ["device"]
