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
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func isValidGitHubWebhookSignature(secret, signature string, payload []byte) bool {
	payloadHash := hmacSha256Hex(secret, payload)
	return hmac.Equal(
		[]byte(payloadHash),
		[]byte(signature),
	)
}

func hmacSha256Hex(secret string, playload []byte) string {
	hm := hmac.New(sha256.New, []byte(secret))
	hm.Write(playload)
	sum := hm.Sum(nil) // nil means "create a new array for me"
	return fmt.Sprintf("sha256=%x", sum)
}
