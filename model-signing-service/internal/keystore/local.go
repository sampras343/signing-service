package keystore

import (
    "crypto/ecdsa"
    "crypto/x509"
    "encoding/pem"
    "os"
)

type LocalKeyStore struct {
    PrivPath string
    PubPath  string
}

func (l *LocalKeyStore) LoadPrivateKey() (*ecdsa.PrivateKey, error) {
    keyData, err := os.ReadFile(l.PrivPath)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(keyData)
    return x509.ParseECPrivateKey(block.Bytes)
}

func (l *LocalKeyStore) LoadPublicKey() (*ecdsa.PublicKey, error) {
    keyData, err := os.ReadFile(l.PubPath)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(keyData)
    pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return pubInterface.(*ecdsa.PublicKey), nil
}
