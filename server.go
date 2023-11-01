package main

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config Config
	gin    *gin.Engine
}

func NewServer(cfg Config) (*Server, error) {
	gin.DefaultWriter = customGinLogWriter{}
	gin.DefaultErrorWriter = customGinLogWriter{}

	r := gin.New()

	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/health"},
		}),
	)

	r.HandleMethodNotAllowed = true

	s := &Server{cfg, r}
	if err := s.bindEndpoints(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Serve(host string, port uint16) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("Serving HTTP.", "address", addr)
	return s.gin.Run(addr)
}

func (s *Server) bindEndpoints() error {
	s.gin.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	for path, endpoint := range s.config.Endpoints {
		if err := s.bindEndpointMethods(path, endpoint); err != nil {
			return fmt.Errorf("endpoint %q: %w", path, err)
		}
	}
	return nil
}

func (s *Server) bindEndpointMethods(path string, endpoint Endpoint) error {
	methods := endpoint.Methods
	if len(methods) == 0 {
		slog.Info("Using default methods for endpoint.",
			"endpoint", path,
			"methods", defaultMethods)
		methods = defaultMethods
	}
	ghWebhookToken, err := readSecret(endpoint.Auth.GitHubWebhookSecret, endpoint.Auth.GitHubWebhookSecretFile)
	if err != nil {
		return err
	}
	client := &http.Client{}
	if endpoint.NoFollowRedirect {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	if endpoint.InsecureSkipVerifyTLS {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	for _, method := range methods {
		handler := endpointHandler{
			Endpoint:           endpoint,
			method:             method,
			endpointPath:       path,
			client:             client,
			githubWebhookToken: ghWebhookToken,
		}
		s.gin.Handle(method, path, handler.handleRequest)
	}
	return nil
}
