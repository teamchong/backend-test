FROM golang:1.18 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 benthos

WORKDIR /go/src/github.com/teamchong/backend-test/
# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY go.* /go/src/github.com/teamchong/backend-test/
COPY ratelimit/go.* /go/src/github.com/teamchong/backend-test/ratelimit/
RUN go mod download 

# Build
COPY . /go/src/github.com/teamchong/backend-test/
COPY ratelimit /go/src/github.com/teamchong/backend-test/ratelimit/
# Tag timetzdata required for busybox base image:
# https://github.com/benthosdev/benthos/issues/897
RUN make TAGS="timetzdata"

# Pack
FROM busybox AS package

LABEL maintainer="Steven Chong <teamchong+github@pm.me>"

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/src/github.com/teamchong/backend-test .
COPY benthos.yaml /benthos.yaml

USER benthos

EXPOSE 4195

ENTRYPOINT ["/backend-test"]

CMD ["-c", "/benthos.yaml"]