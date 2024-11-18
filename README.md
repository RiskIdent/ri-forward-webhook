<!--
SPDX-FileCopyrightText: 2023 Risk.Ident GmbH <contact@riskident.com>

SPDX-License-Identifier: CC-BY-4.0
-->

# ri-forward-webhook

[![REUSE status](https://api.reuse.software/badge/github.com/RiskIdent/ri-forward-webhook)](https://api.reuse.software/info/github.com/RiskIdent/ri-forward-webhook)

> [!NOTE]
> This project is archived because it is no longer in use
> and will therefore no longer be maintained.

HTTP server that validates and forwards webhooks.

Implements [GitHub webhook validation](https://docs.github.com/en/enterprise-server@3.10/webhooks/using-webhooks/validating-webhook-deliveries#validating-webhook-deliveries)
for SHA-256 based signatures (`X-Hub-Signature-256` header).

## Use case

Exposing Jira's webhook endpoint, but not exposing the entire Jira.

This application acts as a gateway from the public github.com over to our
internal Jira. We expose `ri-forward-webhook`, and with the additional
webhook signature validation we can feel safe on not letting bad actors
sending arbitrary requests to our precious Jira.

In addition, in our infrastructure as code repo,
we use [Cloudflare Tunnels](https://www.cloudflare.com/products/tunnel/)
to expose `ri-forward-webhook`, so that nothing else in our Kubernetes cluster
is accidentally exposed to the Internet.

## Configuration

The service is configured using YAML.

Example:

```yaml
endpoints:
  # Proxy all GET and POST requests to http://google.com
  /foo/bar:
    # Destination URI. Must be set.
    forwardTo: http://google.com
    # Which HTTP methods to forward. Default: [POST]
    methods: [GET, POST]
    # If set, then ri-forward-webhook does not follow any redirections.
    # By default, it will follow up to 10 redirects.
    noFollowRedirect: true
    # If set, then ignores all TLS certificate issues.
    insecureSkipVerifyTls: true

  # Proxy all POST requests to http://google.com,
  # but only if a valid X-Hub-Signature-256 header is supplied
  /foo/doo:
    forwardTo: http://google.com
    auth:
      githubWebhookSecret: mySecretKey

  # Proxy all POST requests to http://google.com,
  # but only if a valid X-Hub-Signature-256 header is supplied
  /foo/moo:
    forwardTo: http://google.com
    auth:
      # Reads secret from file
      githubWebhookSecretFile: /path/to/mysecretkey.txt
```

## Development

Requires Go 1.21 (or later)

```bash
go run .
```

## Building

```bash
podman build . -t ghcr.io/riskident/ri-forward-webhook
```

## Releasing

1. Create a new release on GitHub, with "v" prefix on version: <https://github.com/RiskIdent/ri-forward-webhook/releases/new>

2. Write a small changelog, like so:

   ```markdown
   ## Changes (since v0.3.0)

   - Added some feature. (#123)
   ```

3. Our GitHub Action with goreleaser will build and add artifacts to release

## License

This repository complies with the [REUSE recommendations](https://reuse.software/).

Different licenses are used for different files. In general:

- Go code is licensed under GNU General Public License v3.0 or later ([LICENSES/GPL-3.0-or-later.txt](LICENSES/GPL-3.0-or-later.txt)).
- Documentation licensed under Creative Commons Attribution 4.0 International ([LICENSES/CC-BY-4.0.txt](LICENSES/CC-BY-4.0.txt)).
- Miscellaneous files, e.g `.gitignore`, are licensed under CC0 1.0 Universal ([LICENSES/CC0-1.0.txt](LICENSES/CC0-1.0.txt)).

Please see each file's header or accompanied `.license` file for specifics.
