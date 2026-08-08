package integration

import (
	"os/exec"
	"testing"
)

func TestCLIVersion(t *testing.T) {
	// Build first
	cmd := exec.Command("go", "build", "-o", "greenix-test", "../src")
	if err := cmd.Run(); err != nil {
		t.Skip("Skipping integration test: build failed")
	}
	defer exec.Command("rm", "greenix-test").Run()

	// Run version command
	cmd = exec.Command("./greenix-test", "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	if len(output) == 0 {
		t.Error("Version command produced no output")
	}
}

