package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

var (
	_ = hmac.Equal
	_ = sha256.New
)

func TestSignPayloadWithKey(t *testing.T) {
	tests := []struct {
		name     string
		payload  []byte
		key      []byte
		expected string
	}{
		{
			name:     "empty payload and key",
			payload:  []byte(""),
			key:      []byte(""),
			expected: "b613679a0814d9ec772f95d778c35fc5ff1697c493715653c6c712144292c5ad",
		},
		{
			name:     "non-empty payload, empty key",
			payload:  []byte("hello"),
			key:      []byte(""),
			expected: "4352b26e33fe0d769a8922a6ba29004109f01688e26acc9e6cb347e5a5afc4da",
		},
		{
			name:     "empty payload, non-empty key",
			payload:  []byte(""),
			key:      []byte("secret"),
			expected: "f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			name:     "non-empty payload and key",
			payload:  []byte("hello"),
			key:      []byte("secret"),
			expected: "88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SignPayloadWithKey(tt.payload, tt.key)
			gotHex := hex.EncodeToString(got)
			if gotHex != tt.expected {
				t.Errorf("SignPayloadWithKey() = %s, want %s", gotHex, tt.expected)
			}
		})
	}
}
