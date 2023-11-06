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
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	ForwardTo             *YAMLURL     `yaml:"forwardTo"`
	Methods               []string     `yaml:"methods"`
	Auth                  EndpointAuth `yaml:"auth"`
	NoFollowRedirect      bool         `yaml:"noFollowRedirect"`
	InsecureSkipVerifyTLS bool         `yaml:"insecureSkipVerifyTls"`
}

type EndpointAuth struct {
	GitHubWebhookSecret     string `yaml:"githubWebhookSecret"`
	GitHubWebhookSecretFile string `yaml:"githubWebhookSecretFile"`
}

func readConfig(file string) (Config, error) {
	b, err := os.ReadFile(flags.config)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %q: %w", file, err)
	}
	slog.Info("Read config", "file", file)
	return cfg, nil
}

type YAMLURL url.URL

var _ yaml.Unmarshaler = &YAMLURL{}

func (u *YAMLURL) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	parsed, err := url.Parse(str)
	if err != nil {
		return err
	}
	*u = YAMLURL(*parsed)
	return nil
}

func (u *YAMLURL) AsURL() *url.URL {
	return (*url.URL)(u)
}

func (u *YAMLURL) String() string {
	return u.AsURL().String()
}

func readSecret(value string, filename string) (string, error) {
	if value != "" {
		return value, nil
	}
	if filename == "" {
		return "", nil
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("load secret: %w", err)
	}
	if len(b) == 0 {
		return "", fmt.Errorf("load secret: file was empty: %s", filename)
	}
	if !utf8.Valid(b) {
		return "", fmt.Errorf("load secret: file did not contain valid UTF-8: %s", filename)
	}
	s := strings.TrimSpace(string(b))
	if strings.ContainsAny(s, "\n\r") {
		return "", fmt.Errorf("load secret: file contained multiple lines: %s", filename)
	}
	return s, nil
}
