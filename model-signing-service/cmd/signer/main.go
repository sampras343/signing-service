package main

import (
    "github.com/sampras343/signing-service/model-signing-service/internal/artifact/source"
    "github.com/sampras343/signing-service/model-signing-service/internal/keystore"
    "github.com/sampras343/signing-service/model-signing-service/internal/artifact/sink"
    "github.com/sampras343/signing-service/model-signing-service/internal/crypto"
    "github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func main() {
    src := &source.LocalArtifactSource{}
    keyStore := &keystore.LocalKeyStore{PrivPath: "keys/private.pem", PubPath: "keys/public.pem"}
    priv, _ := keyStore.LoadPrivateKey()
    pub, _ := keyStore.LoadPublicKey()
    signer := &crypto.ECDSASigner{PrivateKey: priv, PublicKey: pub}
    sink := &sink.ZipArtifactSink{}

    svc := service.SigningService{
        ArtifactSource: src,
        KeyStore:       keyStore,
        Signer:         signer,
        Sink:           sink,
    }

    svc.SignAndBundle("input-folder", "output/signed_bundle.zip")
}
