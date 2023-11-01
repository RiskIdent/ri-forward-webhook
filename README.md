# ri-forward-webhook

HTTP server that validates and forwards webhooks.

Implements [GitHub webhook validation](https://docs.github.com/en/enterprise-server@3.10/webhooks/using-webhooks/validating-webhook-deliveries#validating-webhook-deliveries)
for SHA-256 based signatures (`X-Hub-Signature-256` header).

## Use case

Exposing Jira's webhook endpoint, but not exposing the entire Jira.

This application acts as a gateway from the public github.com over to our
internal Jira. We expose `ri-forward-webhook`, and with the additional
webhook signature validation we can feel safe on not letting bad actors
sending arbitrary requests to our precious Jira.

In addition, in the [iac repo](https://github.2rioffice.com/platform/iac),
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
podman build . -t docker-riskident.2rioffice.com/platform/ri-forward-webhook
```
