// internal/cli/cli.go
package cli

import (
	"flag"
	"fmt"
	"log"

	"github.com/sampras343/signing-service/model-signing-service/internal/app"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func RunCLI(args []string) {
	fs := flag.NewFlagSet("cli", flag.ExitOnError)

	inputDir := fs.String("input", "", "Input folder containing artifacts")
	outputZip := fs.String("output", "output/signed_bundle.zip", "Output signed bundle path")
	privKey := fs.String("priv", "keys/private.pem", "ECDSA private key path")
	pubKey := fs.String("pub", "keys/public.pem", "ECDSA public key path")
	bundlePath := fs.String("verify", "", "Signed bundle to verify")

	fs.Parse(args)

	if *bundlePath != "" {
		if err := service.VerifyBundle(*bundlePath); err != nil {
			log.Fatal("❌ Verification failed:", err)
		}
		return
	}

	if *inputDir == "" {
		log.Fatal("❌ Input folder path is required. Use --input <path>")
	}

	cfg := app.Config{
		InputDir:    *inputDir,
		OutputZip:   *outputZip,
		PrivKeyPath: *privKey,
		PubKeyPath:  *pubKey,
	}
	svc, err := app.BuildSigningService(cfg)
	if err != nil {
		log.Fatal("❌ Failed to build service:", err)
	}

	if err := svc.SignAndBundle(*inputDir, *outputZip); err != nil {
		log.Fatal("❌ Signing failed:", err)
	}
	fmt.Println("✅ Signed bundle created at:", *outputZip)
}
