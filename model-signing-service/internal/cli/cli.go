package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/sampras343/signing-service/model-signing-service/internal/app"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func RunCLI(args []string) {
	if len(args) < 1 {
		PrintCLIUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "sign":
		runSign(args[1:])
	case "verify":
		runVerify(args[1:])
	case "help", "-h", "--help":
		PrintCLIUsage()
	default:
		fmt.Printf("❌ Unknown command: %s\n", args[0])
		PrintCLIUsage()
		os.Exit(1)
	}
}

func runSign(args []string) {
	fs := flag.NewFlagSet("sign", flag.ExitOnError)

	inputDir := fs.String("input", "", "Path to input folder with model.onnx, dataset.csv, metadata.json")
	outputZip := fs.String("output", "output/signed_bundle.zip", "Output signed bundle path")
	privKey := fs.String("priv", "keys/private.pem", "ECDSA private key path")
	pubKey := fs.String("pub", "keys/public.pem", "ECDSA public key path")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if *inputDir == "" {
		fs.Usage()
		os.Exit(1)
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

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: signing-service verify <signed_bundle.zip>")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}

	bundlePath := fs.Arg(0)
	if err := service.VerifyBundle(bundlePath); err != nil {
		log.Fatal("❌ Verification failed:", err)
	}
}

