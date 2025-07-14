package version

import (
	"testing"
)

func Test_VersionString(t *testing.T) {
	Version = "v1.0.0"
	Date = "2023-01-01"

	expected := "v1.0.0 (built 2023-01-01 with go1.24.4)"
	result := String()

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
