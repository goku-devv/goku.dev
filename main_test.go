package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The stylesheet URL must change whenever the file does: content, templates
// and static deploy without a rebuild, and a build-time version would leave
// returning visitors with new markup on the cached old CSS.
func TestAssetVersion_FollowsTheFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "style.css")
	if err := os.WriteFile(p, []byte("a{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(p, old, old); err != nil {
		t.Fatal(err)
	}
	before := assetVersion(p)
	now := time.Now()
	if err := os.Chtimes(p, now, now); err != nil {
		t.Fatal(err)
	}
	if after := assetVersion(p); before == "" || before == after {
		t.Errorf("version did not follow the file: before=%q after=%q", before, after)
	}

	lastUpdated = "2026-01-02"
	t.Cleanup(func() { lastUpdated = "" })
	if got := assetVersion(filepath.Join(t.TempDir(), "missing.css")); got != "2026-01-02" {
		t.Errorf("missing file: got %q, want the build date", got)
	}
}
