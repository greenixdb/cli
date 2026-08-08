package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// Show help for specific command
			targetCmd, _, err := rootCmd.Find(args)
			if err != nil {
				color.Red("❌ Unknown command: %s", strings.Join(args, " "))
				return
			}
			targetCmd.Help()
			return
		}
		displayHelp()
	},
}

func displayHelp() {
	color.Cyan("██████╗ ██████╗ ███████╗███████╗███╗   ██╗██╗██╗  ██╗")
	color.Cyan("██╔════╝ ██╔══██╗██╔════╝██╔════╝████╗  ██║██║╚██╗██╔╝")
	color.Cyan("██║  ███╗██████╔╝█████╗  █████╗  ██╔██╗ ██║██║ ╚███╔╝ ")
	color.Cyan("██║   ██║██╔══██╗██╔══╝  ██╔══╝  ██║╚██╗██║██║ ██╔██╗ ")
	color.Cyan("╚██████╔╝██║  ██║███████╗███████╗██║ ╚████║██║██╔╝ ██╗")
	color.Cyan(" ╚═════╝ ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═╝")
	color.Cyan("")
	
	color.Set(color.FgGreen)
	fmt.Println("📦 Greenix Studio CLI - Build, deploy, and manage Greenix projects")
	color.Unset()
	fmt.Println()

	fmt.Println("🔧 Commands:")
	fmt.Println("  login      Sign in to your Greenix account")
	fmt.Println("  logout     Sign out of your Greenix account")
	fmt.Println("  whoami     Show the currently signed-in account")
	fmt.Println("  init       Initialize a new Greenix Studio project")
	fmt.Println("  build      Build Greenix Studio projects for multiple platforms")
	fmt.Println("  version    Print version information")
	fmt.Println("  help       Display help information")
	fmt.Println()

	fmt.Println("🚩 Flags:")
	fmt.Println("  -v, --version   Print version information")
	fmt.Println("  -h, --help      Display help information")
	fmt.Println("      --verbose   Verbose output")
	fmt.Println("      --config    Config file (default is .greenix.yaml)")
	fmt.Println()

	fmt.Println("📚 Examples:")
	color.Set(color.FgYellow)
	fmt.Println("  greenix login")
	fmt.Println("  greenix whoami")
	fmt.Println("  greenix logout")
	fmt.Println("  greenix init my-project")
	fmt.Println("  greenix build --target android")
	fmt.Println("  greenix build --target windows --arch x64")
	fmt.Println("  greenix build --target all")
	fmt.Println("  greenix version")
	color.Unset()

	fmt.Println()
	fmt.Println("🔗 Links:")
	fmt.Println("  📖 Documentation: https://docs.greenix.studio")
	fmt.Println("  🐛 Report issues: https://github.com/greenix-studio/cli/issues")
	fmt.Println("  💬 Discord: https://discord.gg/greenix-studio")
}

