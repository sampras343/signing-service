package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/sampras343/signing-service/model-signing-service/internal/crypto"
	"github.com/sampras343/signing-service/model-signing-service/internal/model"
	"github.com/sampras343/signing-service/model-signing-service/internal/util"
)

func VerifyBundle(bundlePath string) error {
	tmpDir, err := util.UnzipTemp(bundlePath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	manifestPath := filepath.Join(tmpDir, "signed_manifest.json")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest model.SignedManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return err
	}

	// Verify hashes
	check := func(file, expected string) error {
		bytes, _ := os.ReadFile(filepath.Join(tmpDir, file))
		hash := fmt.Sprintf("%x", sha256.Sum256(bytes))
		if hash != expected {
			return fmt.Errorf("%s hash mismatch", file)
		}
		return nil
	}

	if err := check("model.onnx", manifest.ModelHash); err != nil {
		return err
	}
	if err := check("dataset.csv", manifest.DatasetHash); err != nil {
		return err
	}

	pubKey, err := x509.ParsePKIXPublicKey(manifest.PublicKey)
	if err != nil {
		return err
	}
	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("public key is not ECDSA")
	}

	sig := manifest.Signature
	manifest.Signature = nil
	verifyBytes, _ := json.Marshal(manifest)
	manifest.Signature = sig

	signer := &crypto.ECDSASigner{PublicKey: ecdsaKey}
	if !signer.Verify(verifyBytes, sig) {
		return fmt.Errorf("signature invalid")
	}

	fmt.Println("✅ Signature verified. Bundle authentic.")
	return nil
}
