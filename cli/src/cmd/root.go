package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/fatih/color"
)

var (
	cfgFile     string
	verbose     bool
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "greenix",
	Short: "Greenix Studio CLI - Build, deploy, and manage Greenix projects",
	Long: `Greenix Studio CLI is a command-line tool for building, deploying,
and managing Greenix Studio projects across multiple platforms.
Complete documentation is available at https://docs.greenix.studio`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			printVersion()
			return nil
		}
		if len(args) > 0 {
			color.Red("❌ Unknown command: %s", args[0])
			fmt.Println()
			displayHelp()
			return fmt.Errorf("unknown command %q", args[0])
		}
		displayHelp()
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			color.Set(color.FgYellow)
			fmt.Println("🔧 Verbose mode enabled")
			color.Unset()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .greenix.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "print version information")

	// Route -h/--help through the branded help screen.
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c == rootCmd {
			displayHelp()
			return
		}
		c.Println(c.UsageString())
	})

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".greenix")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

