package decryptor

import (
	"crypto/rsa"
	"errors"
)

type Decryptor struct {
	PrivateKey *rsa.PrivateKey
	chunkSize  int
}

func NewDecryptor(pk *rsa.PrivateKey) *Decryptor {
	d := Decryptor{PrivateKey: pk}
	if pk != nil {
		d.chunkSize = pk.Size()
	}
	return &d
}

func (d *Decryptor) DecryptData(encryptedData []byte) ([]byte, error) {
	if d.PrivateKey == nil {
		return nil, errors.New("no private key configured")
	}

	if len(encryptedData)%d.chunkSize != 0 {
		return nil, errors.New("invalid encrypted data length")
	}

	var decryptedData []byte
	for i := 0; i < len(encryptedData); i += d.chunkSize {
		chunk := encryptedData[i : i+d.chunkSize]
		decryptedChunk, err := rsa.DecryptPKCS1v15(nil, d.PrivateKey, chunk)
		if err != nil {
			return nil, err
		}
		decryptedData = append(decryptedData, decryptedChunk...)
	}
	return decryptedData, nil
}
