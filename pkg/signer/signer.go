package signer

import (
	"crypto/hmac"
	"crypto/sha256"
)

func SignPayloadWithKey(payload, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(payload)
	return h.Sum(nil)
}
