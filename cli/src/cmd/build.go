package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type BuildTarget struct {
	OS   string
	Arch string
	Ext  string
}

var (
	outputDir string
	targets   []string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Greenix Studio projects",
	Long: `Build Greenix Studio projects for multiple platforms.
	
Examples:
  greenix build --target android
  greenix build --target ios
  greenix build --target windows --arch x64
  greenix build --target all`,
	Run: func(cmd *cobra.Command, args []string) {
		buildProject()
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "build-output", "output directory")
	buildCmd.Flags().StringSliceVar(&targets, "target", []string{}, "target platforms (android, ios, windows, macos, linux, all)")
	buildCmd.Flags().StringSliceVar(&archs, "arch", []string{}, "architectures (x86, x64, arm64, all)")
}

var archs []string

func buildProject() {
	color.Cyan("🚀 Building Greenix Studio...")
	fmt.Println()

	// Read version from root
	version, buildCode := readVersion()
	fmt.Printf("📌 Version: %s (Build %s)\n", version, buildCode)
	fmt.Println()

	// Determine targets to build
	buildTargets := getBuildTargets()
	if len(buildTargets) == 0 {
		color.Red("❌ No targets specified. Use --target flag.")
		return
	}

	// Build each target
	for _, target := range buildTargets {
		color.Cyan("📦 Building %s...", target)
		buildForTarget(target, version, buildCode)
	}

	color.Green("✅ Build complete! Output in: %s", outputDir)
}

func getBuildTargets() []string {
	if len(targets) == 0 {
		return []string{runtime.GOOS}
	}

	// Expand "all" to all supported platforms
	if len(targets) == 1 && targets[0] == "all" {
		return []string{"android", "ios", "windows", "macos", "linux"}
	}

	return targets
}

func buildForTarget(target, version, buildCode string) {
	// Define architecture mapping
	archMap := map[string][]string{
		"android": {"arm64", "x64"},
		"ios":     {"arm64", "x64"},
		"windows": {"x64", "x86", "arm64"},
		"macos":   {"x64", "arm64"},
		"linux":   {"x64", "x86", "arm64"},
	}

	// Determine which architectures to build
	var buildArchs []string
	if len(archs) == 0 || (len(archs) == 1 && archs[0] == "all") {
		buildArchs = archMap[target]
	} else {
		buildArchs = archs
	}

	for _, arch := range buildArchs {
		outputPath := filepath.Join("..", outputDir, "cli", target, arch)
		os.MkdirAll(outputPath, 0755)

		filename := fmt.Sprintf("greenix-%s-%s", version, buildCode)
		if target == "windows" {
			filename += ".exe"
		}

		fullPath := filepath.Join(outputPath, filename)

		color.Yellow("   Building %s/%s -> %s", target, arch, filename)
		
		// Build command
		cmd := exec.Command("go", "build",
			"-ldflags", fmt.Sprintf("-X main.Version=%s -X main.BuildCode=%s -X main.BuildDate=%s",
				version, buildCode, time.Now().Format("2006-01-02 15:04:05")),
			"-o", fullPath,
			"../src",
		)
		
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", target),
			fmt.Sprintf("GOARCH=%s", arch),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			color.Red("   ❌ Failed: %s", err)
			fmt.Println(string(output))
			continue
		}

		color.Green("   ✅ Built successfully")
		
		// Create symlink without version for convenience
		symlinkPath := filepath.Join(outputPath, "greenix")
		if target == "windows" {
			symlinkPath += ".exe"
		}
		os.Remove(symlinkPath)
		os.Symlink(filename, symlinkPath)
	}
}

func readVersion() (string, string) {
	// Try to read from root version.properties
	content, err := os.ReadFile("../../version.properties")
	if err != nil {
		// Fallback to default
		return "dev", "0"
	}

	lines := strings.Split(string(content), "\n")
	var version, buildCode string

	for _, line := range lines {
		if strings.HasPrefix(line, "VERSION_NAME=") {
			version = strings.TrimPrefix(line, "VERSION_NAME=")
			version = strings.TrimSpace(version)
		}
		if strings.HasPrefix(line, "VERSION_CODE=") {
			buildCode = strings.TrimPrefix(line, "VERSION_CODE=")
			buildCode = strings.TrimSpace(buildCode)
		}
	}

	if version == "" {
		version = "dev"
	}
	if buildCode == "" {
		buildCode = "0"
	}

	return version, buildCode
}

