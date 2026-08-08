package main

import (
	"fmt"
	"os"

	"github.com/greenix-studio/cli/cmd"
)

var (
	Version   = "dev"
	BuildCode = "0"
	BuildDate = "unknown"
	GoVersion = ""
)

func main() {
	// Inject build info into commands
	cmd.Version = Version
	cmd.BuildCode = BuildCode
	cmd.BuildDate = BuildDate
	cmd.GoVersion = GoVersion

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

