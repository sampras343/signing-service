package service

import (
    "encoding/json"
    "os"
	"fmt"
	"crypto/sha256"
	"crypto/x509"
	"path/filepath"
	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/source"
	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/sink"
    "github.com/sampras343/signing-service/model-signing-service/internal/model"
	"github.com/sampras343/signing-service/model-signing-service/internal/keystore"
	"github.com/sampras343/signing-service/model-signing-service/internal/crypto"
)

type SigningService struct {
    ArtifactSource source.ArtifactSource
    KeyStore       keystore.KeyStore
    Signer         crypto.Signer
    Sink           sink.ArtifactSink
}

func (s *SigningService) SignAndBundle(inputDir string, outputBundle string) error {
    bundle, err := s.ArtifactSource.Load(inputDir)
    if err != nil {
        return err
    }

    modelBytes, _ := os.ReadFile(bundle.ModelFilePath)
    datasetBytes, _ := os.ReadFile(bundle.DatasetFilePath)
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(s.Signer.GetPublicKey())
	if err != nil {
		return err
	}	

    manifest := &model.SignedManifest{
        ModelHash:   fmt.Sprintf("%x", sha256.Sum256(modelBytes)),
        DatasetHash: fmt.Sprintf("%x", sha256.Sum256(datasetBytes)),
        Metadata:    bundle.Metadata,
        PublicKey:   pubKeyBytes,
    }

    toSign, _ := json.Marshal(manifest)
    signature, _ := s.Signer.Sign(toSign)
    manifest.Signature = signature

    return s.Sink.Write(outputBundle, manifest, []string{
        bundle.ModelFilePath,
        bundle.DatasetFilePath,
        filepath.Join(inputDir, "metadata.json"),
    })
}
