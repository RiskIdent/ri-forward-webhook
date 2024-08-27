# SPDX-FileCopyrightText: 2023 Risk.Ident GmbH <contact@riskident.com>
#
# SPDX-License-Identifier: CC0-1.0

FROM docker.io/library/golang:1.23.0-alpine AS build

WORKDIR /opt/ri-forward-webhook
COPY go.mod go.sum .
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 go install

# NOTE: When updating here, remember to also update in ./goreleaser.Dockerfile
FROM scratch
COPY --from=build /go/bin/ri-forward-webhook /usr/bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["ri-forward-webhook"]
USER 10000

LABEL \
	org.opencontainers.image.source=https://github.com/RiskIdent/ri-forward-webhook \
	org.opencontainers.image.description="Forwards and validates webhooks" \
	org.opencontainers.image.licenses=GPL-3.0-or-later
