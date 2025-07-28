package app

import (
	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/sink"
	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/source"
	"github.com/sampras343/signing-service/model-signing-service/internal/crypto"
	"github.com/sampras343/signing-service/model-signing-service/internal/keystore"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

type Config struct {
	InputDir    string
	OutputZip   string
	PrivKeyPath string
	PubKeyPath  string
}

func BuildSigningService(cfg Config) (*service.SigningService, error) {
	src := &source.LocalArtifactSource{}
	keyStore := &keystore.LocalKeyStore{
		PrivPath: cfg.PrivKeyPath,
		PubPath:  cfg.PubKeyPath,
	}

	priv, err := keyStore.LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	pub, err := keyStore.LoadPublicKey()
	if err != nil {
		return nil, err
	}

	signer := &crypto.ECDSASigner{PrivateKey: priv, PublicKey: pub}
	sink := &sink.ZipArtifactSink{}

	return &service.SigningService{
		ArtifactSource: src,
		KeyStore:       keyStore,
		Signer:         signer,
		Sink:           sink,
	}, nil
}
