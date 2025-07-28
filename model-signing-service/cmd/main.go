package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/sampras343/signing-service/model-signing-service/internal/app"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
	"github.com/sampras343/signing-service/model-signing-service/internal/util"
)

func main() {
	mode := flag.String("mode", "cli", "Execution mode: cli | api")
	inputDir := flag.String("input", "", "Input folder")
	outputZip := flag.String("output", "output/signed_bundle.zip", "Output zip")
	privKey := flag.String("priv", "keys/private.pem", "Private key")
	pubKey := flag.String("pub", "keys/public.pem", "Public key")
	bundlePath := flag.String("verify", "", "Bundle path to verify")
	flag.Parse()

	if *mode == util.CLI_MODE {
		if *bundlePath != "" {
			if err := service.VerifyBundle(*bundlePath); err != nil {
				log.Fatal(err)
			}
			return
		}

		cfg := app.Config{
			InputDir:    *inputDir,
			OutputZip:   *outputZip,
			PrivKeyPath: *privKey,
			PubKeyPath:  *pubKey,
		}
		svc, err := app.BuildSigningService(cfg)
		if err != nil {
			log.Fatal(err)
		}
		if err := svc.SignAndBundle(*inputDir, *outputZip); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ Signed bundle created at:", *outputZip)
	} else {
		fmt.Println("Invalid mode. Use cli or api.")
		os.Exit(1)
	}
}
