package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Greenix Studio project",
	Long: `Initialize a new Greenix Studio project in the current directory.
This creates a basic project structure and configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("🚀 Initializing Greenix Studio project...")
		
		projectName := "my-greenix-project"
		if len(args) > 0 {
			projectName = args[0]
		}

		// Create project structure
		dirs := []string{
			"studio/android-app",
			"studio/ios-app",
			"studio/windows-app",
			"studio/macos-app",
			"studio/linux-app",
			"cli/src",
			"build-output/studio",
			"build-output/cli",
		}

		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				color.Red("❌ Failed to create %s: %v", dir, err)
				return
			}
			color.Green("✅ Created %s", dir)
		}

		// Create version.properties
		versionContent := `# Greenix Studio Project
VERSION_NAME=1.0.0
VERSION_CODE=1
RELEASE_DATE=2026-08-08
PROJECT_NAME=` + projectName + `
`
		err := os.WriteFile("version.properties", []byte(versionContent), 0644)
		if err != nil {
			color.Red("❌ Failed to create version.properties: %v", err)
			return
		}
		color.Green("✅ Created version.properties")

		// Create README
		readmeContent := `# ` + projectName + `

Greenix Studio project.

## Build Commands

- Android: ` + "`./gradlew assembleRelease`" + `
- iOS: ` + "`xcodebuild -scheme ...`" + `
- Windows: ` + "`dotnet build -c Release`" + `
- macOS: ` + "`xcodebuild -scheme ...`" + `
- Linux: ` + "`dotnet build -c Release`" + `

## Version

Current version: 1.0.0 (Build 1)
`
		err = os.WriteFile("README.md", []byte(readmeContent), 0644)
		if err != nil {
			color.Red("❌ Failed to create README.md: %v", err)
			return
		}
		color.Green("✅ Created README.md")

		color.Green("\n✅ Project initialized successfully!")
		color.Cyan("\n📖 Next steps:")
		fmt.Println("   1. cd " + projectName)
		fmt.Println("   2. Edit version.properties to set your version")
		fmt.Println("   3. Build with: greenix build --target all")
	},
}


