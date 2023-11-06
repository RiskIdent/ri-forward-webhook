// SPDX-FileCopyrightText: 2023 Risk.Ident GmbH <contact@riskident.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

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
