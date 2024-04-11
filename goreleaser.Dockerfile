# SPDX-FileCopyrightText: 2024 Risk.Ident GmbH <contact@riskident.com>
#
# SPDX-License-Identifier: CC0-1.0

FROM docker.io/library/alpine AS certs
RUN apk add --no-cache ca-certificates

# NOTE: When updating here, remember to also update in ./Dockerfile
FROM scratch
COPY ri-forward-webhook /usr/local/bin/
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["ri-forward-webhook"]
USER 10000
LABEL \
	org.opencontainers.image.source=https://github.com/RiskIdent/ri-forward-webhook \
	org.opencontainers.image.description="Forwards and validates webhooks" \
	org.opencontainers.image.licenses=GPL-3.0-or-later
