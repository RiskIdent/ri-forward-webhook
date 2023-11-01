package main

import (
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

var defaultMethods = []string{
	"POST",
}

var flags = struct {
	config string
	host   string
	port   uint16
}{
	config: "ri-forward-webhook.yaml",
	host:   "0.0.0.0",
	port:   8080,
}

func init() {
	pflag.StringVarP(&flags.config, "config", "C", flags.config, "Config file to load")
	pflag.StringVar(&flags.host, "host", flags.host, "Host to bind server to")
	pflag.Uint16Var(&flags.port, "port", flags.port, "Port to bind server to")
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	pflag.Parse()
	if err := run(); err != nil {
		slog.Error("Execution failed.", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := readConfig(flags.config)
	if err != nil {
		return err
	}
	server, err := NewServer(cfg)
	if err != nil {
		return err
	}
	return server.Serve(flags.host, flags.port)
}
