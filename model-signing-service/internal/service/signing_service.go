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

func (s *SigningService) SignAndBundle(inputDir, outputBundle string) error {
    bundle, err := s.ArtifactSource.Load(inputDir)
    if err != nil {
        return err
    }

    type hashResult struct {
        name string
        hash string
        err  error
    }
    results := make(chan hashResult, 2)

    go func() {
        modelBytes, err := os.ReadFile(bundle.ModelFilePath)
        if err != nil {
            results <- hashResult{"model", "", err}
            return
        }
        results <- hashResult{"model", fmt.Sprintf("%x", sha256.Sum256(modelBytes)), nil}
    }()

    go func() {
        datasetBytes, err := os.ReadFile(bundle.DatasetFilePath)
        if err != nil {
            results <- hashResult{"dataset", "", err}
            return
        }
        results <- hashResult{"dataset", fmt.Sprintf("%x", sha256.Sum256(datasetBytes)), nil}
    }()

    var modelHash, datasetHash string
    for i := 0; i < 2; i++ {
        r := <-results
        if r.err != nil {
            return r.err
        }
        if r.name == "model" {
            modelHash = r.hash
        } else {
            datasetHash = r.hash
        }
    }

    pubKeyBytes, err := x509.MarshalPKIXPublicKey(s.Signer.GetPublicKey())
    if err != nil {
        return err
    }

    manifest := &model.SignedManifest{
        ModelHash:   modelHash,
        DatasetHash: datasetHash,
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
