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
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

type customGinLogWriter struct {
}

var _ io.Writer = customGinLogWriter{}

func (customGinLogWriter) Write(b []byte) (int, error) {
	level, b := parseLevelAndMessage(b)
	slog.Log(context.Background(), level, fmt.Sprintf("[GIN] %s", b))
	return len(b), nil
}

func parseLevelAndMessage(b []byte) (slog.Level, []byte) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	b = bytes.ReplaceAll(b, []byte("\n"), []byte("\n\t"))

	level := slog.LevelInfo
	if trimmed := bytes.TrimPrefix(b, []byte("[GIN-debug] ")); len(trimmed) != len(b) {
		level = slog.LevelDebug
		b = trimmed
	} else if trimmed := bytes.TrimPrefix(b, []byte("[GIN] ")); len(trimmed) != len(b) {
		b = trimmed
	}

	if trimmed := bytes.TrimPrefix(b, []byte("[WARNING] ")); len(trimmed) != len(b) {
		level = slog.LevelWarn
		b = trimmed
	}

	return level, b
}
