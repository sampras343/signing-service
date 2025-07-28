package util 

import "fmt"

func PrintGlobalHelp() {
	fmt.Println("Usage: signing-service [--mode cli|api] <command> [flags]\n")
	fmt.Println("Modes:")
	fmt.Println("  cli   Run signing or verification from the command line (default)")
	fmt.Println("  api   Run signing-service as an HTTP API server\n")
	fmt.Println("Commands (CLI mode):")
	fmt.Println("  sign     Create a signed bundle")
	fmt.Println("  verify   Verify a signed bundle\n")
	fmt.Println("Examples:")
	fmt.Println("  signing-service sign --input ./sample --output ./bundle.zip")
	fmt.Println("  signing-service verify ./bundle.zip")
	fmt.Println("  signing-service --mode api")
}