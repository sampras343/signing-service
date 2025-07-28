package cli

import (
	"fmt"
)

func PrintCLIUsage() {
	fmt.Println("CLI Mode Usage:")
	fmt.Println("  signing-service sign   --input ./dir --output ./bundle.zip [--priv key.pem --pub key.pem]")
	fmt.Println("  signing-service verify ./bundle.zip")
}