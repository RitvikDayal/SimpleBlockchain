package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // elliptical curve digital signing algorithm
	PublicKey  []byte
}

func GenerateKeyPair() (ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubKey := privateKey.PublicKey
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	return *privateKey, pubKeyBytes
}

func NewWallet() *Wallet {
	privateKey, publicKey := GenerateKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubKeyHash := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubKeyHash[:])
	if err != nil {
		panic(err)
	}
	publicRipemd160 := RIPEMD160Hasher.Sum(nil)
	return publicRipemd160
}

func CheckSum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:checksumLength]
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := PublicKeyHash(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := CheckSum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)

	return address
}
