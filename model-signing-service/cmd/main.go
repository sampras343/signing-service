// cmd/main.go
package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/sampras343/signing-service/model-signing-service/internal/util"
	"github.com/sampras343/signing-service/model-signing-service/internal/cli"
)

func main() {
	mode := util.CLI_MODE
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--mode") {
			if strings.Contains(arg, "=") {
				mode = strings.SplitN(arg, "=", 2)[1]
			} else if i+1 < len(os.Args) {
				mode = os.Args[i+1]
			}
		}
	}
	if envMode := os.Getenv("MODE"); envMode != "" {
		mode = envMode
	}

	fmt.Printf("Starting signing-service in %s mode\n", mode)

	switch mode {
	case util.CLI_MODE:
		filteredArgs := util.FilterModeArgs(os.Args[1:])
		cli.RunCLI(filteredArgs)
	default:
		fmt.Println("Invalid mode. Use cli or api.")
		os.Exit(1)
	}
}
