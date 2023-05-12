package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"log"

	"github.com/ADTEULAA2023/tacochain/pkg"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	//hexadecimal representation of 0
	version = byte(0x00)
)

type Wallet struct {
	//ecdsa = eliptical curve digital signiture algorithm
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func MakeWallet() *Wallet {
	privateKey, publicKey := pkg.NewECDSAKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	if err != nil {
		log.Panic(err)
	}
	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)

	checksum := Checksum(versionedHash)

	finalHash := append(versionedHash, checksum...)
	address := base58.Encode(finalHash)

	return []byte(address)
}
