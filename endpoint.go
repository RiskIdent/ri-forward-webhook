package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type endpointHandler struct {
	Endpoint
	method       string
	endpointPath string
	client       *http.Client

	githubWebhookToken string
}

func (e endpointHandler) handleRequest(c *gin.Context) {
	b, ok := e.readBody(c)
	if !ok {
		return
	}

	if !e.handleRequestAuth(c, b) {
		return
	}

	resp, err := e.doForwardedReq(c, b)
	if err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	e.forwardResponse(c, resp)
}

func (e endpointHandler) readBody(c *gin.Context) ([]byte, bool) {
	defer c.Request.Body.Close()
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError,
			fmt.Errorf("read request body: %w", err))
		return nil, false
	}
	return b, true
}

func (e endpointHandler) doForwardedReq(c *gin.Context, bodyBytes []byte) (*http.Response, error) {
	req, err := e.newForwardedReq(c, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	slog.Info("Forwarded request",
		"endpoint", e.endpointPath,
		"method", e.method,
		"forwardTo", req.URL.String(),
		"status", resp.StatusCode,
		"requestType", req.Header.Get("Content-Type"),
		"responseType", resp.Header.Get("Content-Type"))

	return resp, err
}

func (e endpointHandler) newForwardedReq(c *gin.Context, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(e.method, e.ForwardTo.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = c.Request.Header.Clone()
	return req, nil
}

func (e endpointHandler) forwardResponse(c *gin.Context, resp *http.Response) {
	c.Status(resp.StatusCode)

	headerWriter := c.Writer.Header()
	for key, values := range resp.Header {
		for _, value := range values {
			headerWriter.Add(key, value)
		}
	}

	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		c.Error(err)
		c.AbortWithError(http.StatusBadGateway, err)
	}
}

func (e endpointHandler) handleRequestAuth(c *gin.Context, bodyBytes []byte) bool {
	if e.githubWebhookToken == "" {
		return true
	}
	defer c.Request.Body.Close()

	signature := c.Request.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("missing header: X-Hub-Signature-256"))
		return false
	}

	if !isValidGitHubWebhookSignature(e.githubWebhookToken, signature, bodyBytes) {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("signature did not match: %s", signature))
		return false
	}

	slog.Debug("Request passed X-Hub-Signature-256 verification",
		"endpoint", e.endpointPath,
		"method", e.method,
		"signature", signature)

	return true
}
