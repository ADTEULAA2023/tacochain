package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"log"
	mathrand "math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_/-+!@#$%^&*"

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func randStringBytes() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[mathrand.Intn(16)]
	}
	return string(b)
}

func NewECDSAKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func EncryptTransactionData(data []byte) (string, string, []byte, error) {
	private := randStringBytes()
	public := randStringBytes()
	block, err := aes.NewCipher([]byte(private + public))
	if err != nil {
		return "", "", nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(data))
	cfb.XORKeyStream(cipherText, data)

	return private, public, []byte(base64.StdEncoding.EncodeToString(cipherText)), nil
}

func DecodeTransactionData(privateKey, publicKey, encryptedData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(privateKey + publicKey))
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(data))
	cfb.XORKeyStream(plainText, data)
	return string(plainText), nil
}
