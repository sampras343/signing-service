package crypto

import (
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/sha256"
)

type ECDSASigner struct {
    PrivateKey *ecdsa.PrivateKey
    PublicKey  *ecdsa.PublicKey
}

func (e *ECDSASigner) Sign(data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    return ecdsa.SignASN1(rand.Reader, e.PrivateKey, hash[:])
}

func (e *ECDSASigner) Verify(data []byte, signature []byte) bool {
    hash := sha256.Sum256(data)
    return ecdsa.VerifyASN1(e.PublicKey, hash[:], signature)
}

func (e *ECDSASigner) GetPublicKey() *ecdsa.PublicKey {
    return e.PublicKey
}