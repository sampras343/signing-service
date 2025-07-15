package crypto

import "crypto/ecdsa"

type Signer interface {
    Sign(data []byte) ([]byte, error)
    Verify(data []byte, signature []byte) bool
	GetPublicKey() *ecdsa.PublicKey 
}
