package version

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_VersionString(t *testing.T) {
	Version = "v1.0.0"
	Date = "2023-01-01"
	goVersion := runtime.Version()

	expected := fmt.Sprintf("v1.0.0 (built 2023-01-01 with %s)", goVersion)
	result := String()

	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
