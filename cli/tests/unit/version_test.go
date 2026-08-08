package unit

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// Test version formatting
	version := "1.2.3"
	buildCode := "123"

	expected := "Greenix Studio CLI v1.2.3 (Build 123)"
	actual := formatVersion(version, buildCode)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func formatVersion(version, buildCode string) string {
	return "Greenix Studio CLI v" + version + " (Build " + buildCode + ")"
}

