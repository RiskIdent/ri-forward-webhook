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
