FROM golang as builder
WORKDIR /go/src/github.com/gocityengineering/elastic-metrics
ADD . ./
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE on
RUN \
  go get && \
  go vet && \
  go test -v ./... && \
  go build

FROM ubuntu:21.04
WORKDIR /app/
RUN groupadd app && useradd -g app app
COPY --from=builder /go/src/github.com/gocityengineering/elastic-metrics/elastic-metrics /usr/local/bin/elastic-metrics
COPY config/config_empty.yaml /etc/elastic-metrics/config.yaml
COPY schema/schema.yaml /etc/elastic-metrics/schema.yaml
USER app
CMD ["elastic-metrics"]
