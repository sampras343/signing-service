package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/sink"
	"github.com/sampras343/signing-service/model-signing-service/internal/artifact/source"
	"github.com/sampras343/signing-service/model-signing-service/internal/crypto"
	"github.com/sampras343/signing-service/model-signing-service/internal/keystore"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func main() {
	inputDir := flag.String("input", "", "Path to input folder containing model, dataset, and metadata")
	outputZip := flag.String("output", "output/signed_bundle.zip", "Path to output signed bundle (.zip)")
	privKeyPath := flag.String("priv", "keys/private.pem", "Path to ECDSA private key")
	pubKeyPath := flag.String("pub", "keys/public.pem", "Path to ECDSA public key")

	flag.Parse()

	if *inputDir == "" {
		log.Fatal("❌ Input folder path is required. Use --input <path>")
	}

	src := &source.LocalArtifactSource{}
	keyStore := &keystore.LocalKeyStore{PrivPath: *privKeyPath, PubPath: *pubKeyPath}

	priv, err := keyStore.LoadPrivateKey()
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	pub, err := keyStore.LoadPublicKey()
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	signer := &crypto.ECDSASigner{PrivateKey: priv, PublicKey: pub}
	sink := &sink.ZipArtifactSink{}

	svc := service.SigningService{
		ArtifactSource: src,
		KeyStore:       keyStore,
		Signer:         signer,
		Sink:           sink,
	}

	err = svc.SignAndBundle(*inputDir, *outputZip)
	if err != nil {
		log.Fatalf("Signing failed: %v", err)
	}

	fmt.Println("✅ Signed bundle created at:", *outputZip)
}
