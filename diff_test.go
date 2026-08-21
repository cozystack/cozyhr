package main

import (
	"os"
	"strings"
	"testing"
)

// TestDiffManifests pins the noise-suppression behavior of diffManifests:
// pairs of manifests that differ only in serialization details (Helm
// "# Source:" comment headers, YAML block-scalar chomping style, long-line
// folding, and embedded-JSON whitespace) must produce no diff output, while
// a genuine value-level change must still be reported - even when it's
// accompanied by the same cosmetic noise.
func TestDiffManifests(t *testing.T) {
	tests := []struct {
		name       string
		wantChange bool
	}{
		{name: "comment_header", wantChange: false},
		{name: "block_scalar_chomping", wantChange: false},
		{name: "long_line_folding", wantChange: false},
		{name: "embedded_json", wantChange: false},
		{name: "real_change", wantChange: true},
		{name: "embedded_json_real_change", wantChange: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := readFixture(t, tt.name+".current.yaml")
			desired := readFixture(t, tt.name+".desired.yaml")

			out, err := diffManifests(current, desired, "cozy-monitoring")
			if err != nil {
				t.Fatalf("diffManifests: %v", err)
			}

			gotChange := strings.Contains(out, "has changed")
			if gotChange != tt.wantChange {
				t.Errorf("diffManifests(%s) reported change=%v, want %v; output:\n%s", tt.name, gotChange, tt.wantChange, out)
			}
		})
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/diff/" + name)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return data
}
