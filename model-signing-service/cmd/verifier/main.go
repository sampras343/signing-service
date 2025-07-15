package main

import (
	"archive/zip"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sampras343/signing-service/model-signing-service/internal/crypto"
	"github.com/sampras343/signing-service/model-signing-service/internal/model"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifier <signed_bundle.zip>")
		return
	}

	bundlePath := os.Args[1]
	tmpDir := "tmp_bundle_extract"
	_ = os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	err := unzip(bundlePath, tmpDir)
	if err != nil {
		panic(err)
	}

	manifestPath := filepath.Join(tmpDir, "signed_manifest.json")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		panic(err)
	}

	var manifest model.SignedManifest
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		panic(err)
	}

	modelPath := filepath.Join(tmpDir, "model.onnx")
	modelBytes, _ := os.ReadFile(modelPath)
	modelHash := fmt.Sprintf("%x", sha256Sum(modelBytes))
	if modelHash != manifest.ModelHash {
		fmt.Println("❌ Model hash mismatch!")
		return
	}

	datasetPath := filepath.Join(tmpDir, "dataset.csv")
	datasetBytes, _ := os.ReadFile(datasetPath)
	datasetHash := fmt.Sprintf("%x", sha256Sum(datasetBytes))
	if datasetHash != manifest.DatasetHash {
		fmt.Println("❌ Dataset hash mismatch!")
		return
	}

	pubKey, err := x509.ParsePKIXPublicKey(manifest.PublicKey)
	if err != nil {
		panic(err)
	}

	ecdsaKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("❌ Public key is not ECDSA!")
		return
	}

	// Temporarily strip signature for re-hash
	sig := manifest.Signature
	manifest.Signature = nil
	verifyBytes, _ := json.Marshal(manifest)
	manifest.Signature = sig // restore

	signer := &crypto.ECDSASigner{PublicKey: ecdsaKey}
	if signer.Verify(verifyBytes, sig) {
		fmt.Println("✅ Signature verified. Bundle is authentic and untampered.")
	} else {
		fmt.Println("❌ Signature invalid.")
	}
}

func sha256Sum(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
