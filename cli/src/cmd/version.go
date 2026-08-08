package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildCode = "0"
	BuildDate = "unknown"
	GoVersion = runtime.Version()
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Aliases: []string{"v", "--version"},
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	color.Cyan("██████╗ ██████╗ ███████╗███████╗███╗   ██╗██╗██╗  ██╗")
	color.Cyan("██╔════╝ ██╔══██╗██╔════╝██╔════╝████╗  ██║██║╚██╗██╔╝")
	color.Cyan("██║  ███╗██████╔╝█████╗  █████╗  ██╔██╗ ██║██║ ╚███╔╝ ")
	color.Cyan("██║   ██║██╔══██╗██╔══╝  ██╔══╝  ██║╚██╗██║██║ ██╔██╗ ")
	color.Cyan("╚██████╔╝██║  ██║███████╗███████╗██║ ╚████║██║██╔╝ ██╗")
	color.Cyan(" ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═╝")
	color.Cyan("")
	
	color.Set(color.FgGreen)
	fmt.Printf("📦 Greenix Studio CLI v%s (Build %s)\n", Version, BuildCode)
	color.Unset()
	
	fmt.Printf("   Build Date: %s\n", BuildDate)
	fmt.Printf("   Go Version: %s\n", GoVersion)
	fmt.Printf("   Platform:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	
	color.Set(color.FgYellow)
	fmt.Println("📖 Documentation: https://docs.greenix.studio")
	fmt.Println("🐛 Report issues: https://github.com/greenix-studio/cli/issues")
	color.Unset()
}
