package decryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecryptor(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	t.Run("successful decryption", func(t *testing.T) {
		d := NewDecryptor(privateKey)

		testData := []byte("test data to encrypt")

		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, testData)
		require.NoError(t, err)

		decrypted, err := d.DecryptData(encrypted)
		require.NoError(t, err)
		require.Equal(t, testData, decrypted)
	})

	t.Run("invalid data length", func(t *testing.T) {
		d := NewDecryptor(privateKey)

		invalidData := make([]byte, d.chunkSize+1)

		_, err := d.DecryptData(invalidData)
		require.Error(t, err)
		require.Equal(t, "invalid encrypted data length", err.Error())
	})

	t.Run("decryption failure", func(t *testing.T) {
		d := NewDecryptor(privateKey)

		invalidEncrypted := make([]byte, d.chunkSize)

		_, err := d.DecryptData(invalidEncrypted)
		require.Error(t, err)
	})
}

func TestNewDecryptor(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	d := NewDecryptor(privateKey)
	require.NotNil(t, d)
	require.Equal(t, privateKey, d.PrivateKey)
	require.Equal(t, privateKey.Size(), d.chunkSize)
}
