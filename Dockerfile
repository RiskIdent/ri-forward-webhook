FROM docker.io/library/golang:1.21-alpine AS build

WORKDIR /opt/ri-forward-webhook
COPY go.mod go.sum .
RUN go mod download

COPY *.go .
ENV CGO_ENABLED=0
RUN go install

FROM scratch
COPY --from=build /go/bin/ri-forward-webhook /usr/bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["ri-forward-webhook"]
USER 10000
