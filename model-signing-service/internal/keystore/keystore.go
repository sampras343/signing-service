package keystore

import "crypto/ecdsa"

type KeyStore interface {
    LoadPrivateKey() (*ecdsa.PrivateKey, error)
    LoadPublicKey() (*ecdsa.PublicKey, error)
}
